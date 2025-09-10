package helper

import (
	"fmt"
	"os"

	"github.com/alfianX/jserver/pkg/iso"
	"github.com/moov-io/iso8583"
)

func CreateIsoTest() {
	isomessage := iso8583.NewMessage(iso.Spec87Hex)
	isomessage.MTI("0210")

	isomessage.Field(3, "401000")
	isomessage.Field(4, "000001000000")
	isomessage.Field(11, "000007")
	isomessage.Field(14, "110142")
	isomessage.Field(22, "0708")
	isomessage.Field(24, "017")
	isomessage.Field(25, "000000000007")
	isomessage.Field(35, "00")
	isomessage.Field(41, "71100001")
	isomessage.Field(42, "711000100010004")
	isomessage.Field(52, "910AA767149BA4F533BF3030")
	isomessage.Field(62, "000007")
	isomessage.Field(63, "9900000000151001001001                   MUHAMMAD WAHID FALAN PURY     TARIK TUNAI     00000000019123456123456                SEFTIANI RIDHA WAHYUNI        ")
	isomessage.Field(64, "0000000000000000")

	rawMessage, err := isomessage.Pack()
	if err != nil {
		fmt.Println(err)
	}

	// isomessageUpack := iso8583.NewMessage(Spec87)

	println(string(rawMessage))

	isomessage.Unpack(rawMessage)

	iso8583.Describe(isomessage, os.Stdout)
}
