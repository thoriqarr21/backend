package models

import (
	"time"

	// "gorm.io/gorm"
)

// ─── t_contacts ──────────────────────────────────────────────────────────────

type Contact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"column:first_name;size:100" json:"first_name"`
	LastName  string    `gorm:"column:last_name;size:100" json:"last_name"`
	Phone     string    `gorm:"column:phone;size:30" json:"phone"`
	Email     string    `gorm:"column:email;size:100" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Contact) TableName() string { return "t_contacts" }

// ─── t_flight_transactions ───────────────────────────────────────────────────

type FlightTransaction struct {
	ID                   uint            `gorm:"primaryKey" json:"id"`
	ContactID            uint            `gorm:"column:contact_id" json:"contact_id"`
	Contact              *Contact        `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
	VMCode               string          `gorm:"column:vm_code;size:20" json:"vm_code"`
	PaymentMethodID      uint            `gorm:"column:payment_method_id" json:"payment_method_id"`
	PaymentMethod        *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Amount               float64         `gorm:"column:amount;type:numeric(12,2)" json:"amount"`
	AdminFee             float64         `gorm:"column:admin_fee;type:numeric(12,2);default:0" json:"admin_fee"`
	Markup               float64         `gorm:"column:markup;type:numeric(12,2);default:0" json:"markup"`
	PaymentGatewayCharge float64         `gorm:"column:payment_gateway_charge;type:numeric(12,2);default:0" json:"payment_gateway_charge"`
	TripType             string          `gorm:"column:trip_type;size:20" json:"trip_type"`             // oneway, roundtrip
	Status               string          `gorm:"column:status;size:20;default:'pending'" json:"status"` // pending, completed, failed, cancelled
	SourcePlatform       string          `gorm:"column:source_platform;size:20" json:"source_platform"` // vm, web, mobile
	CreatedAt            time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time       `gorm:"column:updated_at" json:"updated_at"`
	Bookings             []FlightBooking `gorm:"foreignKey:TransactionID" json:"bookings,omitempty"`
}

func (FlightTransaction) TableName() string { return "t_flight_transactions" }

// ─── t_flight_bookings ───────────────────────────────────────────────────────

type FlightBooking struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	TransactionID uint              `gorm:"column:transaction_id" json:"transaction_id"`
	PNRID         string            `gorm:"column:pnrid;size:20" json:"pnrid"`
	BookingCode   string            `gorm:"column:booking_code;size:50" json:"booking_code"`
	NTSA          float64           `gorm:"column:ntsa;type:numeric(12,2);default:0" json:"ntsa"`
	Commission    float64           `gorm:"column:commission;type:numeric(12,2);default:0" json:"commission"`
	PublishFare   float64           `gorm:"column:publish_fare;type:numeric(12,2);default:0" json:"publish_fare"`
	Status        string            `gorm:"column:status;size:20;default:'pending'" json:"status"` // pending, booked, ticketed, cancelled
	CreatedAt     time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at" json:"updated_at"`
	Passengers    []FlightPassenger `gorm:"foreignKey:BookingID" json:"passengers,omitempty"`
	Segments      []FlightSegment   `gorm:"foreignKey:BookingID" json:"segments,omitempty"`
	Remarks       []FlightRemark    `gorm:"foreignKey:BookingID" json:"remarks,omitempty"`
}

func (FlightBooking) TableName() string { return "t_flight_bookings" }

// ─── t_flight_passengers ─────────────────────────────────────────────────────

type FlightPassenger struct {
	ID                     uint       `gorm:"primaryKey" json:"id"`
	BookingID              uint       `gorm:"column:booking_id" json:"booking_id"`
	Title                  string     `gorm:"column:title;size:10" json:"title"`    // Mr, Mrs, Ms, Mstr, Miss
	FullName               string     `gorm:"column:full_name;size:150" json:"full_name"`
	Type                   string     `gorm:"column:type;size:10" json:"type"`      // ADULT, CHILD, INFANT
	PaxID                  int        `gorm:"column:pax_id" json:"pax_id"`
	DOB                    *time.Time `gorm:"column:dob" json:"dob,omitempty"`
	AdultAssoc             *int       `gorm:"column:adult_assoc" json:"adult_assoc,omitempty"`
	PassportNo             string     `gorm:"column:passport_no;size:30" json:"passport_no,omitempty"`
	PassportNationality    string     `gorm:"column:passport_nationality;size:5" json:"passport_nationality,omitempty"`
	PassportIssuingCountry string     `gorm:"column:passport_issuing_country;size:5" json:"passport_issuing_country,omitempty"`
	PassportExpiry         *time.Time `gorm:"column:passport_expiry" json:"passport_expiry,omitempty"`
	TicketNumber           string     `gorm:"column:ticket_number;size:50" json:"ticket_number,omitempty"`
	CreatedAt              time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (FlightPassenger) TableName() string { return "t_flight_passengers" }

// ─── t_flight_segments ───────────────────────────────────────────────────────

type FlightSegment struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	BookingID      uint   `gorm:"column:booking_id" json:"booking_id"`
	AirlineCode    string `gorm:"column:airlinecode;size:5" json:"airlinecode"`
	FlightNo       string `gorm:"column:flightno;size:10" json:"flightno"`
	FromCity       string `gorm:"column:fromcity;size:5" json:"fromcity"`
	ToCity         string `gorm:"column:tocity;size:5" json:"tocity"`
	ClassCode      string `gorm:"column:classcode;size:5" json:"classcode"`
	Cabin          string `gorm:"column:cabin;size:20" json:"cabin"`       // ECONOMY, BUSINESS, FIRST
	DepartTime     string `gorm:"column:departtime;size:10" json:"departtime"`
	ArriveTime     string `gorm:"column:arrivetime;size:10" json:"arrivetime"`
	DepartDate     string `gorm:"column:departdate;size:15" json:"departdate"` // DD-Mon-YYYY
	ArrivalDate    string `gorm:"column:arrivaldate;size:15" json:"arrivaldate"`
	DurationHour   int    `gorm:"column:durationhour" json:"durationhour"`
	DurationMinute int    `gorm:"column:durationminute" json:"durationminute"`
}

func (FlightSegment) TableName() string { return "t_flight_segments" }

// ─── t_flight_remarks ────────────────────────────────────────────────────────

type FlightRemark struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	BookingID uint   `gorm:"column:booking_id" json:"booking_id"`
	Remark    string `gorm:"column:remark;type:text" json:"remark"`
}

func (FlightRemark) TableName() string { return "t_flight_remarks" }

// ─── t_phone_balance_transactions ────────────────────────────────────────────

type PhoneBalanceTransaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Phone       string    `gorm:"column:phone;size:30" json:"phone"`
	Amount      float64   `gorm:"column:amount;type:numeric(14,2)" json:"amount"`
	Type        string    `gorm:"column:type;size:10" json:"type"` // DEDUCT, TOPUP
	ProcessedBy *uint     `gorm:"column:processed_by" json:"processed_by,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PhoneBalanceTransaction) TableName() string { return "t_phone_balance_transactions" }

// ─── REQUEST / INPUT STRUCTS ─────────────────────────────────────────────────

type FlightAvailabilityRequest struct {
	FromCity   string `json:"fromcity"`
	ToCity     string `json:"tocity"`
	DepartDate string `json:"departdate"` // format: YYYY-MM-DD
	ReturnDate string `json:"returndate"` // format: YYYY-MM-DD, opsional
	Adult      string `json:"adult"`
	Child      string `json:"child"`
	Infant     string `json:"infant"`
	Cabin      string `json:"cabin"` // ECONOMY, BUSINESS, FIRST (opsional)
}

type RetrieveFareRequest struct {
	// Trips berisi daftar trip yang akan di-retrieve fare-nya.
	// Untuk OW = 1 trip type "departure"
	// Untuk RT = 2 trips, type "departure" dan "return"
	Trips         []FareTripRequest `json:"trips" binding:"required,min=1"`
	AdultCount    int               `json:"adult_count"`
	ChildCount    int               `json:"child_count"`
	InfantCount   int               `json:"infant_count"`
	WithInsurance bool              `json:"with_insurance"`
}

// FareTripRequest adalah 1 trip untuk AIRFARE request
type FareTripRequest struct {
	TripType    string               `json:"trip_type"`    // "departure" atau "return"
	AirlineCode string               `json:"airline_code"` // diambil dari hasil search
	TripDate    string               `json:"trip_date"`    // YYYY-MM-DD
	Throughfare bool                 `json:"throughfare"`  // true jika connecting throughfare
	ClassCode   string               `json:"class_code"`   // diisi saat throughfare = true
	Segments    []FareSegmentRequest `json:"segments" binding:"required,min=1"`
}

// FareSegmentRequest adalah data segment penerbangan
type FareSegmentRequest struct {
	FlightNo   string `json:"flight_no"`   // diambil dari hasil search
	ClassCode  string `json:"class_code"`  // diisi saat sumlocal, kosong saat throughfare
	FromCity   string `json:"from_city"`   // IATA code
	ToCity     string `json:"to_city"`     // IATA code
	DepartDate string `json:"depart_date"` // YYYY-MM-DD
	DepartTime string `json:"depart_time"` // HH:MM
	ArriveTime string `json:"arrive_time"` // HH:MM
}

type BookFlightRequest struct {
	AdultCount    int           `json:"adult_count"`
	ChildCount    int           `json:"child_count"`
	InfantCount   int           `json:"infant_count"`
	WithInsurance bool          `json:"with_insurance"`
	ServiceFee    int           `json:"service_fee"`    // opsional, 0 = tidak ada
	Trips         []BookTripReq `json:"trips" binding:"required,min=1"`
	ListPax       []Pax         `json:"listpax" binding:"required,min=1"`
	PaxContact    PaxContact    `json:"paxcontact"`
	AgentContact  AgentContact  `json:"agentcontact"`
}

// BookTripReq adalah 1 trip untuk booking
type BookTripReq struct {
	TripType    string           `json:"trip_type"`    // "departure" atau "return"
	AirlineCode string           `json:"airline_code"` // dari hasil search
	TripDate    string           `json:"trip_date"`    // YYYY-MM-DD
	Throughfare bool             `json:"throughfare"`  // true = throughfare, false = sumlocal
	ClassCode   string           `json:"class_code"`   // classcode throughfare (level trip)
	Segments    []BookSegmentReq `json:"segments" binding:"required,min=1"`
}

// BookSegmentReq adalah data per segment dalam 1 trip
type BookSegmentReq struct {
	FlightNo   string `json:"flight_no"`   // dari hasil search
	ClassCode  string `json:"class_code"`  // untuk sumlocal; kosong untuk throughfare
	FromCity   string `json:"from_city"`
	ToCity     string `json:"to_city"`
	DepartDate string `json:"depart_date"` // YYYY-MM-DD
	DepartTime string `json:"depart_time"` // HH:MM
	ArriveTime string `json:"arrive_time"` // HH:MM
}

// AgentContact adalah info kontak agen untuk booking
type AgentContact struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	Phone        string `json:"phone"` // wajib jika with_insurance = true
}

// Pax adalah data penumpang untuk proses booking
type Pax struct {
	Title      string    `json:"title"`       // Mr, Mrs, Ms, Mstr, Miss
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Type       string    `json:"type"`        // ADULT, CHILD, INFANT
	DOB        *string   `json:"dob"`         // YYYY-MM-DD, wajib untuk CHILD & INFANT
	AdultAssoc *int      `json:"adult_assoc"` // pax_id adult yang diasosiasikan (untuk CHILD/INFANT)
	Passport   *Passport `json:"passport"`    // wajib untuk penerbangan internasional
}

// PaxContact adalah data kontak pemesan
type PaxContact struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Phone1       string  `json:"phone1"`
	Phone2       string  `json:"phone2"`       // wajib jika with_insurance = true
	EmailAddress *string `json:"email_address"`
}

// Passport adalah data paspor penumpang
type Passport struct {
	No             string `json:"no"`
	Nationality    string `json:"nationality"`
	IssuingCountry string `json:"issuing_country"`
	Expiry         string `json:"expiry"` // YYYY-MM-DD
}

type RetrievePNRRequest struct {
	PNR string `json:"pnr" binding:"required"`
}

type TicketingRequest struct {
	BookingID uint `json:"booking_id" binding:"required"`
}

type CancelTicketRequest struct {
	BookingID uint   `json:"booking_id" binding:"required"`
	Reason    string `json:"reason"` // disimpan ke remark lokal, tidak dikirim ke Voltras
}

// AutoTicketRequest untuk membuat virtual account pembayaran
type AutoTicketRequest struct {
	BookingID  uint   `json:"booking_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"` // penerima notif VA
	SendTicket bool   `json:"send_ticket"`                    // kirim tiket ke email setelah bayar
}

// AdvanceRetrieveRequest untuk melihat daftar PNR berdasarkan filter
type AdvanceRetrieveRequest struct {
	Airline  string `json:"airline" binding:"required"` // kode maskapai, contoh: GA
	FromDate string `json:"from_date" binding:"required"` // YYYY-MM-DD
	ToDate   string `json:"to_date" binding:"required"`   // YYYY-MM-DD
	Status   string `json:"status"`                       // ALL, BOOKED, TICKETED, CANCELED
}

// ─── PAYMENT STRUCTS ─────────────────────────────────────────────────────────

type CreateTransactionRequest struct {
	BookingID       uint    `json:"booking_id" binding:"required"`
	PaymentMethodID uint    `json:"payment_method_id" binding:"required"`
	BankCode        string  `json:"bank_code,omitempty"`
	Amount          float64 `json:"amount" binding:"required"`
	AdminFee        float64 `json:"admin_fee"`
	Markup          float64 `json:"markup"`
}

type CheckTransactionRequest struct {
	ExternalID string `json:"external_id" binding:"required"`
}

type SendDocumentsRequest struct {
	BookingID uint   `json:"booking_id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type PrintDocumentsRequest struct {
	BookingID uint   `json:"booking_id" binding:"required"`
	DocType   string `json:"doc_type"` // ITINERARY, TICKET, INVOICE
}

// ─── SEGMENT FARE STRUCTS (untuk keperluan lain) ─────────────────────────────

type SumLocalFare struct {
	Segments []SegmentFare `json:"segments"`
}

type ThroughFare struct {
	ClassCode string        `json:"classcode"`
	Segments  []SegmentFare `json:"segments"`
}

type SegmentFare struct {
	FlightNo       string `json:"flightno"`
	ClassCode      string `json:"classcode"`
	FromCity       string `json:"fromcity"`
	ToCity         string `json:"tocity"`
	DepartDate     string `json:"departdate"`
	DepartTime     string `json:"departtime"`
	ArriveTime     string `json:"arrivetime"`
	ArrivalDate    string `json:"arrivaldate"`
	DurationHour   int    `json:"durationhour"`
	DurationMinute int    `json:"durationminute"`
	Cabin          string `json:"cabin"`
}