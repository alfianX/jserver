package iso

type ISO8583BIN87 struct {
	iso8583Bin87format map[int][]int
}

func NewISO8583BIN87() *ISO8583BIN87 {
	return &ISO8583BIN87{
		iso8583Bin87format: map[int][]int{
			0:   {0, 2, 2},
			1:   {1, 2, 8},   // Bit Map Extended
			2:   {1, 2, 8},   // Primary account number (PAN)
			3:   {0, 0, 19},  // Precessing code
			4:   {0, 2, 6},   // Amount transaction
			5:   {0, 2, 12},  // Amount reconciliation
			6:   {0, 2, 12},  // Amount cardholder billing
			7:   {0, 2, 12},  // Date and time transmission
			8:   {0, 2, 10},  // Amount cardholder billing fee
			9:   {0, 2, 8},   // Conversion rate reconciliation
			10:  {0, 2, 8},   // Conversion rate cardholder billing
			11:  {0, 2, 8},   // Systems trace audit number
			12:  {0, 2, 6},   // Date and time local transaction
			13:  {0, 2, 6},   // Date effective
			14:  {0, 2, 4},   // Date expiration
			15:  {0, 2, 4},   // Date settlement
			16:  {0, 2, 4},   // Date conversion
			17:  {0, 2, 4},   // Date capture
			18:  {0, 2, 4},   // Message error indicator
			19:  {0, 2, 4},   // Country code acquiring institution
			20:  {0, 2, 3},   // Country code primary account number (PAN)
			21:  {1, 0, 28},  // Transaction life cycle identification data
			22:  {0, 2, 3},   // Point of service data code
			23:  {0, 2, 3},   // Card sequence number
			24:  {0, 2, 3},   // Function code
			25:  {0, 2, 3},   // Message reason code
			26:  {0, 2, 2},   // Merchant category code
			27:  {0, 2, 2},   // Point of service capability
			28:  {0, 2, 1},   // Date reconciliation
			29:  {0, 2, 8},   // Reconciliation indicator
			30:  {0, 2, 8},   // Amounts original
			31:  {0, 2, 8},   // Acquirer reference number
			32:  {0, 2, 8},   // Acquiring institution identification code
			33:  {0, 0, 11},  // Forwarding institution identification code
			34:  {0, 0, 11},  // Electronic commerce data
			35:  {1, 0, 28},  // Track 2 data
			36:  {2, 0, 37},  // Track 3 data
			37:  {2, 1, 104}, // Retrieval reference number
			38:  {3, 2, 12},  // Approval code
			39:  {3, 2, 6},   // Action code
			40:  {3, 2, 2},   // Service code
			41:  {3, 2, 3},   // Card acceptor terminal identification
			42:  {3, 2, 8},   // Card acceptor identification code
			43:  {3, 2, 15},  // Card acceptor name/location
			44:  {3, 2, 40},  // Additional response data
			45:  {1, 0, 25},  // Track 1 data
			46:  {1, 0, 76},  // Amounts fees
			47:  {1, 1, 999}, // Additional data national
			48:  {1, 1, 999}, // Additional data private
			49:  {3, 1, 999}, // Verification data
			50:  {3, 2, 3},   // Currency code, settlement
			51:  {3, 2, 3},   // Currency code, cardholder billing
			52:  {3, 2, 3},   // Personal identification number (PIN) data
			53:  {1, 2, 8},   // Security related control information
			54:  {0, 2, 16},  // Amounts additional
			55:  {1, 1, 999}, // Integrated circuit card (ICC) system related data
			56:  {3, 1, 999}, // Original data elements
			57:  {1, 1, 999}, // Authorisation life cycle code
			58:  {1, 1, 999}, // Authorising agent institution identification code
			59:  {3, 1, 999}, // Transport data --AX30 24112021
			60:  {3, 1, 999}, // Reserved for national use
			61:  {3, 1, 999}, // Reserved for national use
			62:  {3, 1, 999}, // Reserved for private use
			63:  {3, 1, 999}, // Reserved for private use
			64:  {3, 1, 999}, // Message authentication code (MAC) field
			65:  {1, 2, 8},   // Bitmap tertiary
			66:  {1, 2, 8},   // Settlement code
			67:  {0, 2, 1},   // Extended payment data
			68:  {0, 2, 2},   // Receiving institution country code
			69:  {0, 2, 3},   // Settlement institution county code
			70:  {0, 2, 3},   // Network management Information code
			71:  {0, 2, 3},   // Message number
			72:  {0, 2, 4},   // Data record
			73:  {0, 2, 4},   // Date action
			74:  {0, 2, 6},   // Credits, number
			75:  {0, 2, 10},  // Credits, reversal number
			76:  {0, 2, 10},  // Debits, number
			77:  {0, 2, 10},  // Debits, reversal number
			78:  {0, 2, 10},  // Transfer number
			79:  {0, 2, 10},  // Transfer, reversal number
			80:  {0, 2, 10},  // Inquiries number
			81:  {0, 2, 10},  // Authorizations, number
			82:  {0, 2, 10},  // Credits, processing fee amount
			83:  {0, 2, 12},  // Credits, transaction fee amount
			84:  {0, 2, 12},  // Debits, processing fee amount
			85:  {0, 2, 12},  // Debits, transaction fee amount
			86:  {0, 2, 12},  // Credits, amount
			87:  {0, 2, 15},  // Credits, reversal amount
			88:  {0, 2, 15},  // Debits, amount
			89:  {0, 2, 15},  // Debits, reversal amount
			90:  {0, 2, 15},  // Original data elements
			91:  {0, 2, 42},  // File update code
			92:  {3, 2, 1},   // File security code
			93:  {0, 2, 2},   // Response indicator
			94:  {0, 2, 5},   // Service indicator
			95:  {3, 2, 7},   // Replacement amounts
			96:  {3, 2, 42},  // Message security code
			97:  {3, 2, 8},   // Amount, net settlement
			98:  {0, 2, 16},  // Payee
			99:  {3, 2, 25},  // Settlement institution identification code
			100: {0, 0, 11},  // Receiving institution identification code
			101: {0, 0, 11},  // File name
			102: {3, 2, 17},  // Account identification 1
			103: {3, 0, 28},  // Account identification 2
			104: {3, 0, 28},  // Transaction description
			105: {3, 0, 100}, // Reserved for ISO use
			106: {3, 1, 999}, // Reserved for ISO use
			107: {3, 1, 999}, // Reserved for ISO use
			108: {3, 1, 999}, // Reserved for ISO use
			109: {3, 1, 999}, // Reserved for ISO use
			110: {3, 1, 999}, // Reserved for ISO use
			111: {3, 1, 999}, // Reserved for private use
			112: {3, 1, 999}, // Reserved for private use
			113: {3, 1, 999}, // Reserved for private use
			114: {0, 0, 11},  // Reserved for national use
			115: {3, 1, 999}, // Reserved for national use
			116: {3, 1, 999}, // Reserved for national use
			117: {3, 1, 999}, // Reserved for national use
			118: {3, 1, 999}, // Reserved for national use
			119: {3, 1, 999},
			120: {3, 0, 999}, // Reserved for national use
			121: {3, 1, 999}, // Reserved for private use
			122: {3, 1, 999}, // Reserved for private use
			123: {3, 1, 999}, // Reserved for national use
			124: {3, 1, 999}, // Reserved for private use
			125: {3, 1, 999}, // Info Text
			126: {3, 1, 50},  // Network management information
			127: {3, 1, 6},   // Issuer trace id
			128: {3, 1, 999}, // Reserved for private use
			129: {0, 1, 8},   // Message authentication code (MAC) field
			130: {4, 1, 0},   // Message authentication code (MAC) field
		},
	}
}

func (i *ISO8583BIN87) BitLenType(bit int) int {
	lenType := i.iso8583Bin87format[bit][1]
	return lenType
}

func (i *ISO8583BIN87) BitDataType(bit int) int {
	dataType := i.iso8583Bin87format[bit][0]
	return dataType
}

func (i *ISO8583BIN87) BitLength(bit int) int {
	lenDefault := i.iso8583Bin87format[bit][2]
	return lenDefault
}