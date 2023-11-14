package yaac_backend

import (
	"fmt"
	"image"
	"image/color"
	"math"

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

func order_points(points gocv.PointVector) gocv.PointVector {
	/*
			initialzie a list of coordinates that will be ordered
		    such that the first entry in the list is the top-left,
		    the second entry is the top-right, the third is the
		    bottom-right, and the fourth is the bottom-left
	*/

	pts := points.ToPoints()
	rect := []image.Point{
		image.Point{0, 0},
		image.Point{0, 0},
		image.Point{0, 0},
		image.Point{0, 0},
	}

	// Get Rectangle Points
	arg_min_sum := math.MaxInt
	arg_max_sum := math.MinInt
	arg_min_dif := math.MaxInt
	arg_max_dif := math.MinInt
	for i := 0; i < points.Size(); i++ {
		pnt := points.At(i)

		// the top-left point will have the smallest sum, whereas
		// the bottom-right point will have the largest sum
		sum := pnt.X + pnt.Y
		if sum < arg_min_sum {
			arg_min_sum = sum
		}
		if sum > arg_max_sum {
			arg_max_sum = sum
		}

		// now, compute the difference between the points, the
		// top-right point will have the smallest difference,
		// whereas the bottom-left will have the largest difference
		dif := pnt.Y - pnt.X
		if dif < arg_min_dif {
			arg_min_dif = dif
		}
		if dif > arg_max_dif {
			arg_max_dif = dif
		}
	}
	rect[0] = pts[arg_min_sum]
	rect[2] = pts[arg_max_sum]
	rect[1] = pts[arg_min_dif]
	rect[3] = pts[arg_max_dif]

	// return the ordered coordinates
	return gocv.NewPointVectorFromPoints(rect)
}

/*
	# now, compute the difference between the points, the
    # top-right point will have the smallest difference,
    # whereas the bottom-left will have the largest difference
    diff = np.diff(pts, axis=1)
    rect[1] = pts[np.argmin(diff)]
    rect[3] = pts[np.argmax(diff)]

    # return the ordered coordinates
    return rect
*/

// ## Find the exact (x,y) coordinates of the biggest contour and crop it out
func four_point_transform(img gocv.Mat, pts gocv.PointVector) gocv.Mat {
	// obtain a consistent order of the points and unpack them
	// individually
	rect := order_points(pts)

	//(tl, tr, br, bl) := rect @LEANDER
	rec_arr := rect.ToPoints()
	tl := rec_arr[0]
	tr := rec_arr[1]
	br := rec_arr[2]
	bl := rec_arr[3]

	// compute the width of the new image, which will be the
	// maximum distance between bottom-right and bottom-left
	// x-coordiates or the top-right and top-left x-coordinates
	widthA := math.Sqrt(math.Pow(float64(br.X-bl.X), 2) + math.Pow(float64(br.Y-bl.Y), 2))
	widthB := math.Sqrt(math.Pow(float64(tr.X-tl.X), 2) + math.Pow(float64(tr.Y-tl.Y), 2))
	maxWidth := max(int(widthA), int(widthB))

	// compute the height of the new image, which will be the
	// maximum distance between the top-right and bottom-right
	// y-coordinates or the top-left and bottom-left y-coordinates
	heightA := math.Sqrt(math.Pow(float64(tr.X-br.X), 2) + math.Pow(float64(tr.Y-br.Y), 2))
	heightB := math.Sqrt(math.Pow(float64(tl.X-bl.X), 2) + math.Pow(float64(tl.Y-bl.Y), 2))
	maxHeight := max(int(heightA), int(heightB))

	/*
		now that we have the dimensions of the new image, construct
		the set of destination points to obtain a "birds eye view",
		(i.e. top-down view) of the image, again specifying points
		in the top-left, top-right, bottom-right, and bottom-left
		order
	*/
	dst := gocv.NewPointVectorFromPoints([]image.Point{
		image.Point{0, 0},
		image.Point{maxWidth - 1, 0},
		image.Point{maxWidth - 1, maxHeight - 1},
		image.Point{0, maxHeight - 1},
	})

	// compute the perspective transform matrix and then apply it
	M := gocv.GetPerspectiveTransform(rect, dst)
	warped := gocv.NewMat()
	gocv.WarpPerspective(img, &warped, M, image.Point{maxWidth, maxHeight})

	// return the warped image
	return warped
}

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
