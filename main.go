package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/bmp"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func getPixels(file io.Reader) ([][]Pixel, image.Image, error) {
	img, err := bmp.Decode(file)

	if err != nil {
		return nil, nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}
	return pixels, img, nil
}

func rgbaToPixel(r, g, b, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func calculateSize(multiplication int) (degree,number int) {
	i := 0.0

	for {
		if (int(math.Pow(2, i)) < multiplication) && ((int)(math.Pow(2, i+1)) >= multiplication) {
			return int(i), int(math.Pow(2, i))
		}else {
			i++
		}
	}
}

func printBits(slice []bool) {
	for i := 0; i < len(slice); i++ {
		if slice[i] {
			fmt.Print(1)
		} else {
			fmt.Print(0)
		}
	}
}

func ConvertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, toBase), nil
}
func makeMessage(lenOfMessage int) []uint {
	message := make([]uint, lenOfMessage) // Пример сообщения
	for i := 0; i < cap(message); i++ {
		if i%2 == 1 {
			message[i] = 0x00
		} else {
			message[i] = 0x01
		}
	}

	return message
}

func insertMessage(img image.Image, pixels [][]Pixel, message []uint) (int, int, [][]Pixel) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	countOfMessage := 0
	placeForMessageSize, pow := calculateSize(width * height) // Размер возможной длины сообщения. В столько штук G-компонент будет записан размер сообщения.
	fmt.Println("Number of bits for message size = ", placeForMessageSize)

	flagOfCR := true                                             // Для единичного вывода на экран
	countOfBitPlace := 0                                         // Счетчик записанных битов placeForMessageSize
	binSizeM, _ := ConvertInt(strconv.Itoa(len(message)), 10, 2) // Двоичная СС пространства для сообщения

	if len(message) > pow-len(message) {
		fmt.Println("Size of message too large")
		return 0, 0, nil
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if countOfBitPlace < len(binSizeM) {
				if placeForMessageSize-(y*x+x+len(binSizeM)) > 0 {
					pixels[y][x].G &= 0xFE
					fmt.Print("0")
				} else {
					if binSizeM[countOfBitPlace] == '1' {
						pixels[y][x].G |= 0x01
						fmt.Print("1")
						countOfBitPlace++
					} else {
						pixels[y][x].G &= 0xFE
						fmt.Print("0")
						countOfBitPlace++
					}
				}
			} else {
				if flagOfCR {
					fmt.Println()
					flagOfCR = false
				}
				if (countOfMessage < len(message)) && countOfMessage <= pow {
					if message[countOfMessage] == 0x01 {
						var tmp = pixels[y][x].G
						pixels[y][x].G |= 0x01
						fmt.Println(tmp, "->", pixels[y][x].G, ", ", 1)
						countOfMessage++
					} else {
						var tmp = pixels[y][x].G
						pixels[y][x].G &= 0xFE
						fmt.Println(tmp, "->", pixels[y][x].G, ", ", 0)
						countOfMessage++
					}
				} else {
					return width, height, pixels
				}
			}
		}
	}

	return width, height, pixels
}

func writeImage(width, height int, pixels [][]Pixel) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	newImg := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cyan := color.RGBA{uint8(pixels[y][x].R), uint8(pixels[y][x].G), uint8(pixels[y][x].B), uint8(pixels[y][x].A)}
			newImg.Set(x, y, cyan)
		}
	}

	return newImg
}

func extractMessage(pixelsRec [][]Pixel, imgRec image.Image) []bool {
	boundsRec := imgRec.Bounds()
	widthRec, heightRec := boundsRec.Max.X, boundsRec.Max.Y
	placeForMessageRec, _ := calculateSize(widthRec * heightRec)
	countOfSymbolRec := 0   // Число бит зашифрованного сообщения
	countOfSizeMessage := 0 // Чисто бит пространства для сообщения

	var lengthOfMessageInt int64

	var lengthOfMessageBin []string

	var containerText []bool // Расшифрованное сообщение

	flagOfCountingValue := true

	for y := 0; y < heightRec; y++ {
		for x := 0; x < widthRec; x++ {
			if countOfSizeMessage < placeForMessageRec {
				if pixelsRec[y][x].G%2 == 1 {
					lengthOfMessageBin = append(lengthOfMessageBin, "1")
					countOfSizeMessage++
				} else {
					lengthOfMessageBin = append(lengthOfMessageBin, "0")
					countOfSizeMessage++
				}
			} else {
				if flagOfCountingValue {
					lengthOfMessageInt, _ = strconv.ParseInt(strings.Join(lengthOfMessageBin, ""), 2, 64)
					fmt.Println("Length: ", strings.Join(lengthOfMessageBin, ""), " -> ", lengthOfMessageInt)
					flagOfCountingValue = false
				}
				if countOfSymbolRec < int(lengthOfMessageInt) {
					if pixelsRec[y][x].G%2 == 0 {
						fmt.Println(pixelsRec[y][x].G, " -> ", "0")
						containerText = append(containerText, false)
						countOfSymbolRec++
					} else {
						fmt.Println(pixelsRec[y][x].G, " -> ", "1")
						containerText = append(containerText, true)
						countOfSymbolRec++
					}
				}
			}
		}
	}

	return containerText
}

func main() {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)

	f, err := os.Open("pic\\norm.bmp")

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	pixels, img, err := getPixels(f)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	lenOfMessage := 255
	message := makeMessage(lenOfMessage)
	fmt.Println("Message = ", message, ", length = ", len(message))
	width, height, pixels := insertMessage(img, pixels, message)

	if pixels == nil {
		fmt.Println("\nError!")
		os.Exit(1)
	}

	newImg := writeImage(width, height, pixels)

	s, err := os.Create("pic\\outimage.bmp")

	if err != nil {
		fmt.Println(err)
	}
	bmp.Encode(s, newImg)
	s.Close()
	fmt.Println("\nDecoding")

	file, err := os.Open("pic\\outimage.bmp")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	pixelsRec, imgRec, err := getPixels(file)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	containerText := extractMessage(pixelsRec, imgRec)
	printBits(containerText)
}
