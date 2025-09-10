package iso

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ISOData struct {
	DE map[int]string
}

func NewISOData() *ISOData {
	return &ISOData{
		DE: make(map[int]string, 130),
	}
}

func (iso *ISOData) SetDE(i int, DEData []byte) {
	iso.DE[i] = string(DEData)
}

func CreateIso(DE map[int]string, TPDU string, MTI string, lenType string) string {
	iso := Build(DE, MTI)
	step := TPDU + iso
	ISO := ""

	LenType := lenType
	if LenType == "0" {
		ISO = step
	} else if LenType == "1" { //HEX
		i := len(step) / 2
		ISO = fmt.Sprintf("%04X%s", i, step)
	} else { //BCD
		i := len(step) / 2
		ISO = fmt.Sprintf("%04d%s", i, step)
	}

	return ISO
}

func Build(DE map[int]string, MTI string) string {
	var newISO strings.Builder
	bitLen := 0

	newISO.WriteString(MTI)
	newDE1 := ""
	for I := 2; I <= 64; I++ {
		if DE[I] != "" {
			newDE1 += "1"
		} else {
			newDE1 += "0"
		}
	}

	newDE2 := ""
	for I := 65; I <= 128; I++ {
		if DE[I] != "" {
			newDE2 += "1"
		} else {
			newDE2 += "0"
		}
	}

	if newDE2 == "0000000000000000000000000000000000000000000000000000000000000000" {
		newDE1 = "0" + newDE1
	} else {
		newDE1 = "1" + newDE1
	}

	intDE1 := fromBinary(newDE1)
	DE1Hex := fmt.Sprintf("%X", intDE1)
	DE1Hex = fmt.Sprintf("%016s", DE1Hex) // Pad-Left
	DE[0] = DE1Hex

	DE2Hex := fmt.Sprintf("%X", fromBinary(newDE2))
	DE2Hex = fmt.Sprintf("%016s", DE2Hex) // Pad-Left
	DE[1] = DE2Hex

	if DE2Hex == "0000000000000000" {
		DE[1] = ""
	}

	for I := 0; I <= 128; I++ {
		if DE[I] != "" {
			if I == 100 || I == 32 {
				//
			}

			lenType := NewISO8583BIN87().BitLenType(I + 1)
			lenSize := NewISO8583BIN87().BitLength(I + 1)

			var li, paddedli, BMPadded, sBM string

			bitLen = 0

			switch NewISO8583BIN87().BitDataType(I + 1) {
			case 1:
				switch lenType {
				case 2:
					bitLen = lenSize * 2

					BMPadded = fmt.Sprintf("%0"+fmt.Sprint(bitLen)+"s", DE[I])
					sBM = BMPadded[:bitLen]
					newISO.WriteString(sBM)
				case 0:
					bitLen = len(DE[I])
					if bitLen%2 != 0 {
						bitLen++
					}
					li = fmt.Sprintf("%02d", bitLen/2)

					BMPadded = fmt.Sprintf("%0"+fmt.Sprint(bitLen)+"s", DE[I])
					sBM = li + BMPadded[:bitLen]
					newISO.WriteString(sBM)
				case 1:
					bitLen = len(DE[I])

					if bitLen%2 != 0 {
						bitLen++
					}

					dataHex := fmt.Sprintf("%0"+fmt.Sprint(bitLen)+"s", DE[I])
					var length int

					if I == 55 && false {
						bitLen = len(dataHex)
						length = bitLen / 2
						dataAscii := hex.EncodeToString([]byte([]byte(dataHex)))
						li = hex.EncodeToString([]byte([]byte(fmt.Sprintf("%03d", length))))
						BMPadded = dataAscii

						sBM = li + BMPadded
						newISO.WriteString(sBM)
					} else {
						length = bitLen / 2

						li = fmt.Sprintf("%04d", length)
						BMPadded = dataHex
						sBM = li + BMPadded[:bitLen]
						newISO.WriteString(sBM)
					}
				default:
				}
			case 3:
				hexs := hex.EncodeToString([]byte([]byte(DE[I])))
				if (I == 127 && false) || I == 55 {
					hexs = DE[I]
				}
				switch lenType {
				case 2:
					if I == 49 && false {
						if lenSize%2 != 0 {
							lenSize++
						}
						BMPadded = fmt.Sprintf("%0"+fmt.Sprint(lenSize)+"s", DE[I])
						sBM = BMPadded[:lenSize]
						newISO.WriteString(sBM)
						break
					}

					BMPadded = fmt.Sprintf("%0"+fmt.Sprint(lenSize)+"s", DE[I])
					BMPadded = hex.EncodeToString([]byte([]byte(BMPadded)))
					// BMPadded = hex + strings.Repeat("30", (lenSize*2)-len(hex))
					sBM = BMPadded[:lenSize*2]
					newISO.WriteString(sBM)
				case 0:
					li = fmt.Sprint(len(DE[I]))
					paddedli = fmt.Sprintf("%02s", li)
					newISO.WriteString(paddedli + hexs)
				case 1:
					if I == 127 && false {
						li = fmt.Sprint(len(DE[I]))
						paddedli = fmt.Sprintf("%02s", li)
					} else {
						if I == 55 {
							li = fmt.Sprint(len(DE[I]) / 2)
						} else {
							li = fmt.Sprint(len(DE[I]))
						}
						paddedli = fmt.Sprintf("%04s", li)
					}
					newISO.WriteString(paddedli + hexs)
				default:
				}
			case 0:
				switch lenType {
				case 2:
					if I == 22 && false {
						bitLen = 6
						BMPadded = hex.EncodeToString([]byte([]byte(DE[22])))
						sBM = BMPadded[:bitLen]
						newISO.WriteString(sBM)
					} else {
						bitLen = lenSize
						if lenSize%2 != 0 {
							bitLen++
						}

						BMPadded = fmt.Sprintf("%0"+fmt.Sprint(bitLen)+"s", DE[I])
						sBM = BMPadded[:bitLen]
						newISO.WriteString(sBM)
					}
				case 0:
					bitLen = len(DE[I])
					if len(DE[I])%2 != 0 {
						// bitlen++
						DE[I] = PadRight(DE[I], bitLen+1, '0')
					}

					li = fmt.Sprint(bitLen)
					paddedli = fmt.Sprintf("%02s", li)
					newISO.WriteString(paddedli + DE[I])
				case 1:
					li = fmt.Sprintf("%04d", len(DE[I]))
					paddedli = li
					newISO.WriteString(paddedli + DE[I])
				default:
				}
			case 2:
				switch lenType {
				case 2:
					bitLen = lenSize
					if lenSize%2 != 0 {
						bitLen++
					}

					BMPadded = fmt.Sprintf("%0"+fmt.Sprint(bitLen)+"s", DE[I])
					sBM = BMPadded[:bitLen]
					newISO.WriteString(sBM)
				case 0:
					bitLen = len(DE[I])
					if len(DE[I])%2 != 0 {
						bitLen++
					}

					li = fmt.Sprint(len(DE[I]))
					paddedli = fmt.Sprintf("%02s", li)
					newISO.WriteString(paddedli + PadRight(DE[I], bitLen, 'F'))
				case 1:
					li = fmt.Sprintf("%04d", len(DE[I]))
					paddedli = li
					newISO.WriteString(paddedli + PadRight(DE[I], bitLen, 'F'))
				default:
				}
			}
		}
	}

	return strings.ToUpper(newISO.String())
}

func Parse(iso string) (map[int]string, error) {
	var err error
	DE := make(map[int]string)

	de1binary := strings.Repeat("0", 64)
	de2binary := strings.Repeat("0", 64)
	fieldNo := 129

	myPos := 0
	myLength := 4

	var MTI string

	if len(iso) >= myLength {
		MTI = iso[myPos : myPos+myLength]
	} else {
		return nil, errors.New("MTI error")
	}

	// ===BM 129 is MTI===
	fieldNo = 129
	DE[fieldNo] = MTI
	// ===================
	myPos += myLength
	myLength = 16

	if len(iso) >= myPos+myLength {
		DE[0] = iso[myPos : myPos+myLength]
	} else {
		return nil, errors.New("bitmap error")
	}
	// Convert BM-0 to binary
	de1binary = HexBin(DE[0])

	fieldNo = 1
	if string(de1binary[fieldNo-1]) == "1" {
		myPos += myLength
		myLength = 16
		if len(iso) >= myPos+myLength {
			DE[fieldNo] = iso[myPos : myPos+myLength]
		} else {
			return nil, errors.New("secondary bitmap error")
		}
		de2binary = HexBin(DE[fieldNo])
	}

	var lenType, dataType, lenDefault, lenData int
	// var field, val string

	for fieldNo := 2; fieldNo <= 128; fieldNo++ {
		if fieldNo == 100 || fieldNo == 62 {
			//
		}

		if fieldNo <= 64 {
			if string(de1binary[fieldNo-1]) != "1" {
				continue
			}
		} else {
			//===CHECK POINT FOR BM# 65 to 128============
			// if de2binary == "" {
			//     return DE
			// }
			if de2binary == strings.Repeat("0", 64) {
				return DE, nil
			}
			//=============================================
			if string(de2binary[fieldNo-65]) != "1" {
				continue
			}
		}

		myPos += myLength
		lenType = NewISO8583BIN87().BitLenType(fieldNo + 1)
		lenDefault = NewISO8583BIN87().BitLength(fieldNo + 1)
		dataType = NewISO8583BIN87().BitDataType(fieldNo + 1)
		switch lenType {
		case 2:
			lenData = lenDefault
		case 0:
			var ll string
			if len(iso) >= myPos+2 {
				ll = iso[myPos : myPos+2]
			} else {
				return nil, errors.New("LLV " + strconv.Itoa(fieldNo) + " error")
			}
			lenData, err = strconv.Atoi(ll)
			if err != nil {
				return nil, errors.New("LLV " + strconv.Itoa(fieldNo) + " error")
			}
			myPos += 2
		case 1:
			var lll string
			if len(iso) >= myPos+4 {
				lll = iso[myPos+1 : myPos+4]
			} else {
				return nil, errors.New("LLLV " + strconv.Itoa(fieldNo) + " error")
			}
			lenData, err = strconv.Atoi(lll)
			if err != nil {
				return nil, errors.New("LLLV " + strconv.Itoa(fieldNo) + " error")
			}
			myPos += 4
		}
		switch dataType {
		case 3:
			var data string
			if len(iso) >= myPos {
				hexs := iso[myPos:]

				myLength = lenData * 2
				if len(hexs) >= myLength {
					if fieldNo == 55 {
						data = string(hexs[:myLength])
					} else {
						str, err := hex.DecodeString(hexs[:myLength])
						if err != nil {
							return nil, errors.New("DE-" + strconv.Itoa(fieldNo) + " decode hex error")
						}
						data = string(str)
					}
					DE[fieldNo] = data
				} else {
					return nil, errors.New("DE-" + strconv.Itoa(fieldNo) + " error")
				}
			} else {
				return nil, errors.New("DE-" + strconv.Itoa(fieldNo) + " error")
			}
		case 1:
			myLength = lenData * 2

			if len(iso) >= myPos+myLength {
				DE[fieldNo] = iso[myPos : myPos+myLength]
			} else {
				return nil, errors.New("DE-" + strconv.Itoa(fieldNo) + " error")
			}

			if myLength%2 != 0 {
				myPos += 1
			}
		case 0, 2:
			myLength = lenData
			if len(iso) >= myPos+myLength {
				if fieldNo == 35 || fieldNo == 2 {
					if myLength%2 != 0 {
						myLength++
					}
					DE[fieldNo] = iso[myPos : myPos+lenData]
				} else if myLength%2 != 0 {
					if fieldNo == 100 || fieldNo == 32 || fieldNo == 33 || fieldNo == 49 {
						// lenData++
						myLength = lenData + 1
					} else {
						myPos += 1
					}
				}
				DE[fieldNo] = iso[myPos : myPos+lenData]
			} else {
				return nil, errors.New("DE-" + strconv.Itoa(fieldNo) + " error")
			}
		}
	}

	return DE, nil
}

func ParseUp(iso string) map[int]string {
	DE := make(map[int]string)

	de1binary := strings.Repeat("0", 64)
	de2binary := strings.Repeat("0", 64)
	fieldNo := 129

	myPos := 0
	myLength := 4

	MTI := iso[myPos : myPos+myLength]

	// ===BM 129 is MTI===
	fieldNo = 129
	DE[fieldNo] = MTI
	// ===================
	myPos += myLength
	myLength = 16
	DE[0] = iso[myPos : myPos+myLength]

	// Convert BM-0 to binary
	de1binary = HexBin(DE[0])

	fieldNo = 1
	if string(de1binary[fieldNo-1]) == "1" {
		myPos += myLength
		myLength = 16
		DE[fieldNo] = iso[myPos : myPos+myLength]
		de2binary = HexBin(DE[fieldNo])
	}

	var lenType, dataType, lenDefault, lenData int

	for fieldNo := 2; fieldNo <= 22; fieldNo++ {
		if fieldNo == 100 || fieldNo == 62 {
			//
		}

		if fieldNo <= 64 {
			if string(de1binary[fieldNo-1]) != "1" {
				continue
			}
		} else {
			//===CHECK POINT FOR BM# 65 to 128============
			// if de2binary == "" {
			//     return DE
			// }
			if de2binary == strings.Repeat("0", 64) {
				return DE
			}
			//=============================================
			if string(de2binary[fieldNo-65]) != "1" {
				continue
			}
		}

		myPos += myLength
		lenType = NewISO8583BIN87().BitLenType(fieldNo + 1)
		lenDefault = NewISO8583BIN87().BitLength(fieldNo + 1)
		dataType = NewISO8583BIN87().BitDataType(fieldNo + 1)
		switch lenType {
		case 2:
			lenData = lenDefault
		case 0:
			lenData = int(iso[myPos])<<8 + int(iso[myPos+1])
			myPos += 2
		case 1:
			lenData = int(iso[myPos+1])<<8 + int(iso[myPos+2])<<16 + int(iso[myPos+3])<<24
			myPos += 4
		}
		switch dataType {
		case 3:
			hex := iso[myPos:]
			myLength = lenData * 2
			DE[fieldNo] = hex[:myLength]
		case 1:
			myLength = lenData * 2
			DE[fieldNo] = iso[myPos : myPos+myLength]
			if myLength%2 != 0 {
				myPos += 1
			}
		case 0, 2:
			myLength = lenData
			if fieldNo == 35 || fieldNo == 2 {
				if myLength%2 != 0 {
					myLength++
				}
				DE[fieldNo] = iso[myPos : myPos+lenData]
			} else if myLength%2 != 0 {
				if fieldNo == 100 || fieldNo == 32 || fieldNo == 33 || fieldNo == 49 {
					lenData++
					myLength = lenData
				} else {
					myPos += 1
				}
			}
			DE[fieldNo] = iso[myPos : myPos+lenData]
		}
	}

	return DE
}

func DEtoBinary(HexDE string) string {
	deBinary := ""
	for i := 0; i <= 15; i++ {
		deBinary += HexBin(string(HexDE[i]))
	}
	return deBinary
}

func HexBin(sHex string) string {
	sReturn := ""
	sHex = strings.ToUpper(sHex)
	for i := 0; i < len(sHex); i++ {
		switch string(sHex[i]) {
		case "0":
			sReturn += "0000"
		case "1":
			sReturn += "0001"
		case "2":
			sReturn += "0010"
		case "3":
			sReturn += "0011"
		case "4":
			sReturn += "0100"
		case "5":
			sReturn += "0101"
		case "6":
			sReturn += "0110"
		case "7":
			sReturn += "0111"
		case "8":
			sReturn += "1000"
		case "9":
			sReturn += "1001"
		case "A":
			sReturn += "1010"
		case "B":
			sReturn += "1011"
		case "C":
			sReturn += "1100"
		case "D":
			sReturn += "1101"
		case "E":
			sReturn += "1110"
		case "F":
			sReturn += "1111"
		}
	}
	return sReturn
}

func fromBinary(binary string) uint64 {
	var value uint64
	for _, c := range binary {
		value = value*2 + uint64(c-'0')
	}
	return value
}

func Bin2hex(s string) string {
	var hexStr string
	for i := 0; i < len(s); i += 8 {
		end := i + 8
		if end > len(s) {
			end = len(s)
		}
		block := s[i:end]
		value := fromBinary(block)
		hexStr += fmt.Sprintf("%02X", value)
	}
	return hexStr
}

func PadRight(s string, length int, padChar byte) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}
