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
		blurred_threshold := transformation(image)
		image = final_image(blurred_threshold)
		fmt.Println("Done!")
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
func biggest_contour(contours gocv.PointsVector, min_area float64) (int, gocv.PointVector) {
	max_area := 0.0
	biggest_n := 0
	var approx_contour gocv.PointVector

	for n := 0; n < contours.Size(); n++ {
		cnt := contours.At(n)
		area := gocv.ContourArea(cnt)

		if area > float64(min_area) {
			peri := gocv.ArcLength(cnt, true)
			approx := gocv.ApproxPolyDP(cnt, 0.02*peri, true)
			if area > max_area && approx.Size() == 4 {
				max_area = area
				biggest_n = n
				approx_contour = approx
			}
		}
	}

	return biggest_n, approx_contour
}

func order_points(pts gocv.PointVector) gocv.PointVector {
	// ToDo
}

// ## Find the exact (x,y) coordinates of the biggest contour and crop it out
func four_point_transform(image gocv.Mat, pts gocv.PointVector) gocv.Mat {
	// obtain a consistent order of the points and unpack them
    // individually
    rect := order_points(pts)
    (tl, tr, br, bl) := rect
}

/*
def four_point_transform(image, pts):
    # obtain a consistent order of the points and unpack them
    # individually
    rect = order_points(pts)
    (tl, tr, br, bl) = rect

    # compute the width of the new image, which will be the
    # maximum distance between bottom-right and bottom-left
    # x-coordiates or the top-right and top-left x-coordinates
    widthA = np.sqrt(((br[0] - bl[0]) ** 2) + ((br[1] - bl[1]) ** 2))
    widthB = np.sqrt(((tr[0] - tl[0]) ** 2) + ((tr[1] - tl[1]) ** 2))
    maxWidth = max(int(widthA), int(widthB))

    # compute the height of the new image, which will be the
    # maximum distance between the top-right and bottom-right
    # y-coordinates or the top-left and bottom-left y-coordinates
    heightA = np.sqrt(((tr[0] - br[0]) ** 2) + ((tr[1] - br[1]) ** 2))
    heightB = np.sqrt(((tl[0] - bl[0]) ** 2) + ((tl[1] - bl[1]) ** 2))
    maxHeight = max(int(heightA), int(heightB))

    # now that we have the dimensions of the new image, construct
    # the set of destination points to obtain a "birds eye view",
    # (i.e. top-down view) of the image, again specifying points
    # in the top-left, top-right, bottom-right, and bottom-left
    # order
    dst = np.array([
        [0, 0],
        [maxWidth - 1, 0],
        [maxWidth - 1, maxHeight - 1],
        [0, maxHeight - 1]], dtype="float32")

    # compute the perspective transform matrix and then apply it
    M = cv2.getPerspectiveTransform(rect, dst)
    warped = cv2.warpPerspective(image, M, (maxWidth, maxHeight))

    # return the warped image
    return warped
 */

/*
# Transformation the image
**1. Convert the image to grayscale**
**2. Remove noise and smoothen out the image by applying blurring and thresholding techniques**
**3. Use Canny Edge Detection to find the edges**
**4. Find the biggest contour and crop it out**
*/
func transformation(im gocv.Mat) gocv.Mat {
	// threshold == image == gray?

	image := gocv.NewMat()
	im.CopyTo(&image)

	/*
		height := image.Size()[0]
		width := image.Size()[0]
		channels := image.Channels()
	*/

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

		if poly.Size() != 4 {
			continue
		}

		simplified_contours.Append(poly)
	}

	_, approx_contour := biggest_contour(simplified_contours, float64(image_size[0]*image_size[1]))

	gocv.DrawContours(&image, simplified_contours, -1, color.RGBA{0, 0, 0, 0}, 1)

	var dst = gocv.NewMat()
	if approx_contour.IsNil() == false && approx_contour.Size() == 4 {
		// approx_contour = np.float32(approx_contour)
		dst = four_point_transform(threshold, approx_contour)
	}
	return dst
}

// **Sharpen the image using Kernel Sharpening Technique**
func final_image(rotated gocv.Mat) gocv.Mat {
	// Create our shapening kernel, it must equal to one eventually
	kernel_sharpening := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV32S) // @LEANDER
	kernel_sharpening.SetIntAt(0, 0, 0)
	kernel_sharpening.SetIntAt(0, 1, -1)
	kernel_sharpening.SetIntAt(0, 2, 0)
	kernel_sharpening.SetIntAt(1, 0, -1)
	kernel_sharpening.SetIntAt(1, 1, 5)
	kernel_sharpening.SetIntAt(1, 2, -1)
	kernel_sharpening.SetIntAt(2, 0, 0)
	kernel_sharpening.SetIntAt(2, 1, -1)
	kernel_sharpening.SetIntAt(2, 2, 0)
	// applying the sharpening kernel to the input image & displaying it.
	sharpened := gocv.NewMat()
	gocv.Filter2D(rotated, &sharpened, -1, kernel_sharpening, image.Point{-1, -1}, 0, gocv.BorderDefault)

	return sharpened
}
