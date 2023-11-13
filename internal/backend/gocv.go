package yaac_backend

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func (b *Backend) StartGoCV() {
	fmt.Println("CV")
	image := gocv.IMRead("./assets/list_3.jpg", gocv.IMReadColor)
	if image.Empty() {
		fmt.Println("Empty image")
	} else {
		fmt.Println("Loaded image")
		transformation(image)
	}
}

func blur_and_threshold(gray gocv.Mat) gocv.Mat {
	gocv.GaussianBlur(gray, &gray, image.Point{3, 3}, 2, 0, gocv.BorderDefault)
	threshold := gocv.Mat{}
	gocv.AdaptiveThreshold(gray, &threshold, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 11, 2)
	gocv.FastNlMeansDenoisingWithParams(threshold, &threshold, 11, 31, 9)
	return threshold
}

// ## **Find the Biggest Contour**
// **Note: We made sure the minimum contour is bigger than 1/10 size of the whole picture. This helps in removing very small contours (noise) from our dataset**
func biggest_contour(contours gocv.PointsVector, min_area []int) (any, any) {
	// ToDo
}

func order_points(pts any) {
	// ToDo
}

// ## Find the exact (x,y) coordinates of the biggest contour and crop it out
func four_point_transform(image any, pts any) {
	// ToDo
}

/*
# Transformation the image
**1. Convert the image to grayscale**
**2. Remove noise and smoothen out the image by applying blurring and thresholding techniques**
**3. Use Canny Edge Detection to find the edges**
**4. Find the biggest contour and crop it out**
*/

func transformation(im gocv.Mat) {
	image := gocv.NewMat()
	im.CopyTo(&image)

	height := image.Size()[0]
	width := image.Size()[0]
	channels := image.Channels()

	gray := gocv.NewMat()
	gocv.CvtColor(image, &gray, gocv.ColorBGRToGray)
	image_size := gray.Size()

	threshold := blur_and_threshold(gray)
	/*
	   # We need two threshold values, minVal and maxVal. Any edges with intensity gradient more than maxVal
	   # are sure to be edges and those below minVal are sure to be non-edges, so discarded.
	   #  Those who lie between these two thresholds are classified edges or non-edges based on their connectivity.
	   # If they are connected to "sure-edge" pixels, they are considered to be part of edges.
	   #  Otherwise, they are also discarded
	*/
	edges := gocv.NewMat()
	gocv.Canny(threshold, &edges, 50, 150) // FIXME - apertureSize=7 UNKNOWN
	hierarchy := gocv.NewMat()
	contours := gocv.FindContoursWithParams(edges, &hierarchy, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var simplified_contours gocv.PointsVector

	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		hull := gocv.NewMat()
		gocv.ConvexHull(cnt, &hull, false, true)
		poly := gocv.ApproxPolyDP(gocv.NewPointVectorFromMat(hull), (0.001 * gocv.ArcLength(gocv.NewPointVectorFromMat(hull), true)), true)
		/*
					Dafuq ist das fürn if statement???
					if poly.shape[0] != 4:
			           continue
		*/
		simplified_contours.Append(poly)
	}

	// simplified_contours = np.array(simplified_contours) - Not neccecary in go?
	biggest_n, approx_contour := biggest_contour(simplified_contours, image_size)

	gocv.DrawContours(&image, simplified_contours, -1, color.RGBA{0, 0, 0, 0}, 1)
}

/*
   simplified_contours = np.array(simplified_contours)
   biggest_n, approx_contour = biggest_contour(simplified_contours, image_size)

   threshold = cv2.drawContours(image, simplified_contours, -1, (0, 0, 0), 1)

   dst = 0
   if approx_contour is not None and len(approx_contour) == 4:
       approx_contour = np.float32(approx_contour)
       dst = four_point_transform(threshold, approx_contour)
   croppedImage = dst
   return croppedImage
*/
