// package voltras

// import (
// 	"crypto/sha1"
// 	"context"
// 	"encoding/hex"
// 	"encoding/xml"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"time"

// 	"mobile-api/models"
// )

// // ─── SIGNATURE ───────────────────────────────────────────────────────────────

// func generateSignature(message, token string) string {
// 	h := sha1.New()
// 	h.Write([]byte(message + token))
// 	return hex.EncodeToString(h.Sum(nil))
// }

// // ─── CLIENT ──────────────────────────────────────────────────────────────────

// type Client struct {
// 	URL        string
// 	OfficeCode string
// 	Token      string
// 	HTTPClient *http.Client
// }

// func NewClient(url, officeCode, token string) *Client {
// 	return &Client{
// 		URL:        url,
// 		OfficeCode: officeCode,
// 		Token:      token,
// 		HTTPClient: &http.Client{Timeout: 90 * time.Second},
// 	}
// }

// // ─── POST REQUEST ────────────────────────────────────────────────────────────

// func (c *Client) post(message string) ([]byte, error) {
// 	eMessage := generateSignature(message, c.Token)

// 	form := url.Values{}
// 	form.Set("OFFICE_CODE", c.OfficeCode)
// 	form.Set("MESSAGE", message)
// 	form.Set("E_MESSAGE", eMessage)

// 	// ⬇️ TAMBAHKAN INI
// 	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
// 	defer cancel()

// 	req, err := http.NewRequestWithContext(ctx, "POST", c.URL, strings.NewReader(form.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	start := time.Now() // debug

// 	resp, err := c.HTTPClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	raw, _ := io.ReadAll(resp.Body)

// 	fmt.Println("=== VOLTRAS REQUEST ===")
// 	fmt.Println(message)
// 	fmt.Println("=== VOLTRAS RESPONSE STATUS:", resp.Status, "===")
// 	fmt.Println("Duration:", time.Since(start)) // ⬅️ penting buat debug
// 	fmt.Println(string(raw))
// 	fmt.Println("=== END ===")

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, string(raw))
// 	}

// 	return raw, nil
// }

// // ─── XML HELPER ──────────────────────────────────────────────────────────────

// func parseXML(data []byte, v interface{}) error {
// 	if err := xml.Unmarshal(data, v); err != nil {
// 		return fmt.Errorf("xml parse error: %w", err)
// 	}
// 	return nil
// }

// // ─── OFFICE INFO ─────────────────────────────────────────────────────────────

// type OfficeInfoResponse struct {
// 	OfficeCode string  `xml:"Body>GetOfficeInformationResponse>OfficeCode"`
// 	OfficeName string  `xml:"Body>GetOfficeInformationResponse>OfficeName"`
// 	Balance    float64 `xml:"Body>GetOfficeInformationResponse>Balance"`
// 	Currency   string  `xml:"Body>GetOfficeInformationResponse>Currency"`
// }

// func (c *Client) GetOfficeInformation() (*OfficeInfoResponse, error) {
// 	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
// <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
//   <Body>
//     <GetOfficeInformation>
//       <OfficeCode>%s</OfficeCode>
//       <Token>%s</Token>
//     </GetOfficeInformation>
//   </Body>
// </Envelope>`, c.OfficeCode, c.Token)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res OfficeInfoResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// // ─── FLIGHT AVAILABILITY ─────────────────────────────────────────────────────
// //
// // Response XML structure dari Voltras:
// //
// // <response type="FLIGHTAVAILABILITY">
// //   <trip type="departure">
// //     <listitem>
// //       <item>
// //         <durationHour>...</durationHour>        ← camelCase!
// //         <durationMinute>...</durationMinute>    ← camelCase!
// //         <throughfare>...</throughfare>
// //         <dayexchange>...</dayexchange>
// //         <listflight>
// //           <flight>
// //             <flightno>...</flightno>
// //             <airlinecode>...</airlinecode>
// //             <fromcity>...</fromcity>
// //             <tocity>...</tocity>
// //             <departdate>...</departdate>
// //             <arrivaldate>...</arrivaldate>
// //             <departtime>...</departtime>
// //             <arrivetime>...</arrivetime>
// //             <durationHour>...</durationHour>    ← camelCase!
// //             <durationMinute>...</durationMinute>← camelCase!
// //             <dayexchange>...</dayexchange>
// //             <meal>...</meal>
// //           </flight>
// //         </listflight>
// //         <listclassgroup>
// //           <classgroup>
// //             <class>
// //               <code>...</code>
// //               <availcode>...</availcode>
// //               <fare>...</fare>
// //               <currency>...</currency>
// //               <pairid>...</pairid>
// //               <baggagekilos>...</baggagekilos>
// //               <cabin>...</cabin>
// //               <servicefee>...</servicefee>
// //             </class>
// //           </classgroup>
// //         </listclassgroup>
// //       </item>
// //     </listitem>
// //   </trip>
// // </response>

// type VolClass struct {
// 	Code         string  `xml:"code"`
// 	AvailCode    string  `xml:"availcode"`
// 	Fare         float64 `xml:"fare"`
// 	Currency     string  `xml:"currency"`
// 	PairID       string  `xml:"pairid"`
// 	BaggageKilos string  `xml:"baggagekilos"`
// 	Cabin        string  `xml:"cabin"`
// 	ServiceFee   float64 `xml:"servicefee"`
// }

// type VolClassGroup struct {
// 	Classes []VolClass `xml:"class"`
// }

// type VolFlight struct {
// 	FlightNo       string `xml:"flightno"`
// 	AirlineCode    string `xml:"airlinecode"`
// 	FromCity       string `xml:"fromcity"`
// 	ToCity         string `xml:"tocity"`
// 	DepartDate     string `xml:"departdate"`
// 	ArrivalDate    string `xml:"arrivaldate"`
// 	DepartTime     string `xml:"departtime"`
// 	ArriveTime     string `xml:"arrivetime"`
// 	DurationHour   int    `xml:"durationHour"`   // camelCase sesuai response XML
// 	DurationMinute int    `xml:"durationMinute"` // camelCase sesuai response XML
// 	DayExchange    int    `xml:"dayexchange"`
// 	Meal           bool   `xml:"meal"`
// }

// type VolFlightItem struct {
// 	DurationHour   int    `xml:"durationHour"`   // camelCase
// 	DurationMinute int    `xml:"durationMinute"` // camelCase
// 	Throughfare    bool   `xml:"throughfare"`
// 	DayExchange    int    `xml:"dayexchange"`
// 	ListFlight     struct {
// 		Flights []VolFlight `xml:"flight"`
// 	} `xml:"listflight"`
// 	ListClassGroup []VolClassGroup `xml:"listclassgroup>classgroup"`
// }

// type VolTrip struct {
// 	Type     string `xml:"type,attr"` // "departure" atau "return"
// 	ListItem struct {
// 		Items []VolFlightItem `xml:"item"`
// 	} `xml:"listitem"`
// }

// // Root: <response type="FLIGHTAVAILABILITY">
// type SearchFlightResponse struct {
// 	XMLName  xml.Name  `xml:"response"`
// 	Trips    []VolTrip `xml:"trip"`
// 	ErrorMsg string    `xml:"description"`
// }

// func (c *Client) FlightAvailability(origin, dest, date string, adt, chd, inf int) (*SearchFlightResponse, error) {
// 	// Parse & validasi tanggal
// 	parsed, err := time.Parse("2006-01-02", date)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid date format (gunakan YYYY-MM-DD): %v", err)
// 	}

// 	loc, _ := time.LoadLocation("Asia/Jakarta")
// 	parsed = parsed.In(loc)
// 	today := time.Now().In(loc).Truncate(24 * time.Hour)

// 	if parsed.Before(today) {
// 		return nil, fmt.Errorf("tanggal keberangkatan tidak boleh kurang dari hari ini")
// 	}

// 	formattedDate := parsed.Format("02-Jan-2006")

// 	body := fmt.Sprintf(`<request type="FLIGHTAVAILABILITY">
// 	<fromcity>%s</fromcity>
// 	<tocity>%s</tocity>
// 	<departdate>%s</departdate>
// 	<listairline>
// 		<airline>GA</airline>
// 		<airline>JT</airline>
// 	</listairline>
// 	<cabin></cabin>
// 	<adult>%d</adult>
// 	<child>%d</child>
// 	<infant>%d</infant>
// 	<cheapestclass>true</cheapestclass>
// </request>`, origin, dest, formattedDate, adt, chd, inf)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res SearchFlightResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}

// 	return &res, nil
// }

// // ─── RETRIEVE FARE (AIRFARE) ─────────────────────────────────────────────────
// //
// // Request: <request type="AIRFARE"> — bukan SOAP!
// // Response:
// //   <response type="AIRFARE">
// //     <totalfare>
// //       <airline code="GA">7752500</airline>
// //     </totalfare>
// //     <servicefee>
// //       <airline code="GA">775250</airline>
// //     </servicefee>
// //     <insurance>0</insurance>
// //   </response>

// // AirlineFare mewakili <airline code="GA">7752500</airline>
// type AirlineFare struct {
// 	Code   string  `xml:"code,attr"`
// 	Amount float64 `xml:",chardata"`
// }

// // AirfareResponse adalah response dari <request type="AIRFARE">
// type AirfareResponse struct {
// 	XMLName    xml.Name      `xml:"response"`
// 	TotalFares []AirlineFare `xml:"totalfare>airline"`
// 	ServiceFees []AirlineFare `xml:"servicefee>airline"`
// 	Insurance  float64       `xml:"insurance"`
// 	ErrorCode  string        `xml:"code"`
// 	ErrorMsg   string        `xml:"description"`
// }

// // RetrieveFareSegment adalah segment penerbangan untuk AIRFARE request
// type RetrieveFareSegment struct {
// 	FlightNo   string
// 	ClassCode  string
// 	FromCity   string
// 	ToCity     string
// 	DepartDate string // format: DD-Mon-YYYY (sudah diformat sebelum masuk)
// 	DepartTime string
// 	ArriveTime string
// }

// // RetrieveFareTrip adalah 1 trip (departure / return) untuk AIRFARE request
// type RetrieveFareTrip struct {
// 	TripType    string // "departure" atau "return"
// 	AirlineCode string
// 	TripDate    string // format: DD-Mon-YYYY
// 	Throughfare bool   // true = throughfare, false = sumlocal
// 	Segments    []RetrieveFareSegment
// 	ClassCode   string // dipakai saat throughfare = true
// }

// func (c *Client) RetrieveFare(
// 	trips []RetrieveFareTrip,
// 	adultCount, childCount, infantCount int,
// 	withInsurance bool,
// ) (*AirfareResponse, error) {

// 	var sb strings.Builder
// 	sb.WriteString(`<request type="AIRFARE">`)

// 	for _, t := range trips {
// 		sb.WriteString(`<trip type="` + t.TripType + `">`)
// 		sb.WriteString(`<airlinecode>` + t.AirlineCode + `</airlinecode>`)
// 		sb.WriteString(`<tripdate>` + t.TripDate + `</tripdate>`)
// 		sb.WriteString(`<listselecteditem>`)
// 		sb.WriteString(`<selecteditem>`)

// 		if t.Throughfare {
// 			// ── THROUGHFARE ──
// 			sb.WriteString(`<throughfare>`)
// 			sb.WriteString(`<listflight>`)
// 			for _, seg := range t.Segments {
// 				sb.WriteString(`<segment>`)
// 				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
// 				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
// 				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
// 				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
// 				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
// 				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
// 				sb.WriteString(`</segment>`)
// 			}
// 			sb.WriteString(`</listflight>`)
// 			sb.WriteString(`<classcode>` + t.ClassCode + `</classcode>`)
// 			sb.WriteString(`</throughfare>`)
// 		} else {
// 			// ── SUM LOCAL (direct atau connecting sum local) ──
// 			sb.WriteString(`<sumlocal>`)
// 			for _, seg := range t.Segments {
// 				sb.WriteString(`<segment>`)
// 				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
// 				sb.WriteString(`<classcode>` + seg.ClassCode + `</classcode>`)
// 				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
// 				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
// 				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
// 				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
// 				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
// 				sb.WriteString(`</segment>`)
// 			}
// 			sb.WriteString(`</sumlocal>`)
// 		}

// 		sb.WriteString(`</selecteditem>`)
// 		sb.WriteString(`</listselecteditem>`)
// 		sb.WriteString(`</trip>`)
// 	}

// 	ins := "false"
// 	if withInsurance {
// 		ins = "true"
// 	}
// 	sb.WriteString(`<withinsurance>` + ins + `</withinsurance>`)
// 	sb.WriteString(fmt.Sprintf(`<adultcount>%d</adultcount>`, adultCount))
// 	sb.WriteString(fmt.Sprintf(`<childcount>%d</childcount>`, childCount))
// 	sb.WriteString(fmt.Sprintf(`<infantcount>%d</infantcount>`, infantCount))
// 	sb.WriteString(`</request>`)

// 	raw, err := c.post(sb.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res AirfareResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}

// 	// Cek error dari Voltras
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}

// 	return &res, nil
// }

// // ─── BOOK FLIGHT ─────────────────────────────────────────────────────────────
// //
// // Request: <request type="BOOK">
// // Response:
// //   <response type="BOOK">
// //     <departure>
// //       <pnrid>14773</pnrid>
// //       <ntsa>1689467</ntsa>
// //       <commission>19633</commission>
// //       <totalfare>1709100</totalfare>
// //       <servicefee>0</servicefee>
// //     </departure>
// //     <return>                    ← hanya muncul jika return dengan maskapai BERBEDA
// //       <pnrid>...</pnrid>
// //       ...
// //     </return>
// //   </response>

// type BookPassenger struct {
// 	Type           string
// 	Title          string
// 	FirstName      string
// 	LastName       string
// 	ID             string // KTP
// 	DOB            string // format: dd-MMM-yyyy (sudah diformat sebelum masuk)
// 	AdultAssoc     string // pax_id adult untuk INFANT
// 	Nationality    string
// 	PassportNo     string
// 	PassportExp    string // format: dd-MMM-yyyy
// 	IssuingCountry string
// }

// // BookTripSegment adalah segment per trip untuk booking
// type BookTripSegment struct {
// 	FlightNo   string
// 	ClassCode  string // dipakai untuk sumlocal; kosong untuk throughfare
// 	FromCity   string
// 	ToCity     string
// 	DepartDate string // format: DD-Mon-YYYY (sudah diformat)
// 	DepartTime string
// 	ArriveTime string
// }

// // BookTrip adalah 1 trip (departure/return) untuk booking
// type BookTrip struct {
// 	TripType    string // "departure" atau "return"
// 	AirlineCode string
// 	TripDate    string            // format: DD-Mon-YYYY
// 	Throughfare bool              // true = throughfare, false = sumlocal
// 	ClassCode   string            // classcode untuk throughfare (level trip)
// 	Segments    []BookTripSegment
// }

// // BookFlightPNR adalah hasil booking per trip
// type BookFlightPNR struct {
// 	PNRID      string  `xml:"pnrid"`
// 	NTSA       float64 `xml:"ntsa"`
// 	Commission float64 `xml:"commission"`
// 	TotalFare  float64 `xml:"totalfare"`
// 	ServiceFee float64 `xml:"servicefee"`
// }

// // BookFlightResponse adalah root response dari Voltras booking
// // <response type="BOOK"><departure>...</departure><return>...</return></response>
// type BookFlightResponse struct {
// 	XMLName   xml.Name      `xml:"response"`
// 	Departure BookFlightPNR `xml:"departure"`
// 	Return    BookFlightPNR `xml:"return"`    // kosong jika OW atau RT maskapai sama
// 	ErrorCode string        `xml:"code"`
// 	ErrorMsg  string        `xml:"description"`
// }

// func (c *Client) BookFlight(
// 	trips []BookTrip,
// 	adultCount, childCount, infantCount int,
// 	withInsurance bool,
// 	serviceFee int,
// 	contact models.PaxContact,
// 	agentContact models.AgentContact,
// 	passengers []BookPassenger,
// ) (*BookFlightResponse, error) {

// 	var sb strings.Builder

// 	sb.WriteString(`<request type="BOOK">`)
// 	sb.WriteString(`<selectionrequest>`)

// 	// withinsurance wajib ada di dalam selectionrequest (sebelum trip)
// 	ins := "false"
// 	if withInsurance {
// 		ins = "true"
// 	}
// 	sb.WriteString(`<withinsurance>` + ins + `</withinsurance>`)
// 	sb.WriteString(fmt.Sprintf(`<adultcount>%d</adultcount>`, adultCount))
// 	sb.WriteString(fmt.Sprintf(`<childcount>%d</childcount>`, childCount))
// 	sb.WriteString(fmt.Sprintf(`<infantcount>%d</infantcount>`, infantCount))

// 	// ── TRIPS ──
// 	for _, t := range trips {
// 		sb.WriteString(`<trip type="` + t.TripType + `">`)
// 		sb.WriteString(`<airlinecode>` + t.AirlineCode + `</airlinecode>`)
// 		sb.WriteString(`<tripdate>` + t.TripDate + `</tripdate>`)
// 		sb.WriteString(`<listselecteditem>`)
// 		sb.WriteString(`<selecteditem>`)

// 		if t.Throughfare {
// 			// ── THROUGHFARE: classcode di level trip, segment tanpa classcode ──
// 			sb.WriteString(`<throughfare>`)
// 			for _, seg := range t.Segments {
// 				sb.WriteString(`<segment>`)
// 				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
// 				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
// 				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
// 				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
// 				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
// 				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
// 				sb.WriteString(`</segment>`)
// 			}
// 			sb.WriteString(`<classcode>` + t.ClassCode + `</classcode>`)
// 			sb.WriteString(`</throughfare>`)
// 		} else {
// 			// ── SUMLOCAL: classcode di tiap segment (direct atau connecting) ──
// 			sb.WriteString(`<sumlocal>`)
// 			for _, seg := range t.Segments {
// 				sb.WriteString(`<segment>`)
// 				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
// 				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
// 				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
// 				sb.WriteString(`<classcode>` + seg.ClassCode + `</classcode>`)
// 				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
// 				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
// 				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
// 				sb.WriteString(`</segment>`)
// 			}
// 			sb.WriteString(`</sumlocal>`)
// 		}

// 		sb.WriteString(`</selecteditem>`)
// 		sb.WriteString(`</listselecteditem>`)
// 		sb.WriteString(`</trip>`)
// 	}

// 	sb.WriteString(`</selectionrequest>`)

// 	// ── LISTPAX ──
// 	sb.WriteString(`<listpax>`)
// 	for _, p := range passengers {
// 		sb.WriteString(`<pax>`)
// 		sb.WriteString(`<title>` + p.Title + `</title>`)
// 		sb.WriteString(`<firstname>` + p.FirstName + `</firstname>`)
// 		sb.WriteString(`<lastname>` + p.LastName + `</lastname>`)
// 		sb.WriteString(`<type>` + p.Type + `</type>`)

// 		// id (KTP) — selalu emit tag meski kosong (sesuai contoh Voltras)
// 		sb.WriteString(`<id>` + p.ID + `</id>`)

// 		// dob — selalu emit tag meski kosong
// 		sb.WriteString(`<dob>` + p.DOB + `</dob>`)

// 		// adultassoc — untuk INFANT wajib
// 		sb.WriteString(`<adultassoc>` + p.AdultAssoc + `</adultassoc>`)

// 		// passport — selalu emit (sesuai contoh Voltras, isi kosong untuk domestik)
// 		sb.WriteString(`<passport>`)
// 		sb.WriteString(`<no>` + p.PassportNo + `</no>`)
// 		sb.WriteString(`<nationality>` + p.Nationality + `</nationality>`)
// 		sb.WriteString(`<issuingcountry>` + p.IssuingCountry + `</issuingcountry>`)
// 		sb.WriteString(`<expiry>` + p.PassportExp + `</expiry>`)
// 		sb.WriteString(`</passport>`)

// 		sb.WriteString(`</pax>`)
// 	}
// 	sb.WriteString(`</listpax>`)

// 	// ── PAXCONTACT ──
// 	sb.WriteString(`<paxcontact>`)
// 	sb.WriteString(`<firstname>` + contact.FirstName + `</firstname>`)
// 	sb.WriteString(`<lastname>` + contact.LastName + `</lastname>`)
// 	sb.WriteString(`<phone1>` + contact.Phone1 + `</phone1>`)
// 	if contact.Phone2 != "" {
// 		sb.WriteString(`<phone2>` + contact.Phone2 + `</phone2>`)
// 	}
// 	if contact.EmailAddress != nil && *contact.EmailAddress != "" {
// 		sb.WriteString(`<emailaddress>` + *contact.EmailAddress + `</emailaddress>`)
// 	}
// 	sb.WriteString(`</paxcontact>`)

// 	// ── AGENTCONTACT ──
// 	sb.WriteString(`<agentcontact>`)
// 	sb.WriteString(`<firstname>` + agentContact.FirstName + `</firstname>`)
// 	sb.WriteString(`<lastname>` + agentContact.LastName + `</lastname>`)
// 	sb.WriteString(`<emailaddress>` + agentContact.EmailAddress + `</emailaddress>`)
// 	if agentContact.Phone != "" {
// 		sb.WriteString(`<phone>` + agentContact.Phone + `</phone>`)
// 	}
// 	sb.WriteString(`</agentcontact>`)

// 	// ── SERVICE FEE (opsional) ──
// 	if serviceFee > 0 {
// 		sb.WriteString(fmt.Sprintf(`<servicefee>%d</servicefee>`, serviceFee))
// 	}

// 	sb.WriteString(`</request>`)

// 	raw, err := c.post(sb.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res BookFlightResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}

// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}

// 	return &res, nil
// }

// // ─── RETRIEVE PNR ────────────────────────────────────────────────────────────
// //
// // Request:  <request type="RETRIEVE"><pnrid>...</pnrid></request>
// // Response: <response type="RETRIEVE"><bookingcode>...</bookingcode>...<listpax>...<listflight>...

// type RetrievePaxDetail struct {
// 	PaxNo      string `xml:"paxno"`
// 	Title      string `xml:"title"`
// 	Name       string `xml:"name"`
// 	TicketNo   string `xml:"ticketno"`
// 	AdultAssoc string `xml:"adultassoc"`
// 	DOB        string `xml:"dob"`
// }

// type RetrieveFlightDetail struct {
// 	AirlineCode string `xml:"airlinecode"`
// 	FlightNo    string `xml:"flightno"`
// 	Departure   string `xml:"departure"`  // kota asal
// 	Arrival     string `xml:"arrival"`    // kota tujuan
// 	DepartDate  string `xml:"departdate"`
// 	DepartTime  string `xml:"departtime"`
// 	ArrivalTime string `xml:"arrivaltime"`
// 	ClassCode   string `xml:"classcode"`
// 	Type        string `xml:"type"` // GO or BACK
// }

// type RetrievePaxContact struct {
// 	FirstName string `xml:"firstname"`
// 	LastName  string `xml:"lastname"`
// 	Phone1    string `xml:"phone1"`
// 	Phone2    string `xml:"phone2"`
// 	Email     string `xml:"email"`
// }

// type RetrieveNTSA struct {
// 	Fare      float64 `xml:"fare"`
// 	Insurance float64 `xml:"insurance"`
// }

// type RetrievePublish struct {
// 	Fare       float64 `xml:"fare"`
// 	Insurance  float64 `xml:"insurance"`
// 	ServiceFee float64 `xml:"serviceFee"`
// }

// type RetrieveAutoTicket struct {
// 	No        string `xml:"no"`
// 	Amount    string `xml:"amount"`
// 	TimeLimit string `xml:"timelimit"`
// 	Status    string `xml:"status"`
// }

// // RetrievePNRResponse adalah response dari <request type="RETRIEVE">
// type RetrievePNRResponse struct {
// 	XMLName     xml.Name               `xml:"response"`
// 	BookingCode string                 `xml:"bookingcode"`
// 	AirlineCode string                 `xml:"airlinecode"`
// 	Status      string                 `xml:"status"`    // BOOKED, TICKETED, CANCELED
// 	TimeLimit   string                 `xml:"timelimit"` // tidak null jika BOOKED
// 	NTSA        RetrieveNTSA           `xml:"ntsa"`
// 	Publish     RetrievePublish        `xml:"publish"`
// 	ListPax     []RetrievePaxDetail    `xml:"listpax>pax"`
// 	PaxContact  RetrievePaxContact     `xml:"paxcontact"`
// 	ListFlight  []RetrieveFlightDetail `xml:"listflight>flight"`
// 	ListRemark  []string               `xml:"listremark>remark"`
// 	AutoTicket  RetrieveAutoTicket     `xml:"autoticket"`
// 	ErrorCode   string                 `xml:"code"`
// 	ErrorMsg    string                 `xml:"description"`
// }

// func (c *Client) RetrievePNR(pnr string) (*RetrievePNRResponse, error) {
// 	body := fmt.Sprintf(`<request type="RETRIEVE"><pnrid>%s</pnrid></request>`, pnr)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res RetrievePNRResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}
// 	return &res, nil
// }

// // ─── TICKETING ───────────────────────────────────────────────────────────────
// //
// // Request:  <request type="TICKET"><pnrid>...</pnrid></request>
// // Response: <response type="TICKET"><pnrid>...</pnrid></response>

// type TicketingResponse struct {
// 	XMLName   xml.Name `xml:"response"`
// 	PNRID     string   `xml:"pnrid"`
// 	ErrorCode string   `xml:"code"`
// 	ErrorMsg  string   `xml:"description"`
// }

// func (c *Client) TicketingFlight(pnr string) (*TicketingResponse, error) {
// 	body := fmt.Sprintf(`<request type="TICKET"><pnrid>%s</pnrid></request>`, pnr)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res TicketingResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}
// 	return &res, nil
// }

// // ─── AUTO TICKET ─────────────────────────────────────────────────────────────
// //
// // Request:  <request type="AUTOTICKET"><pnrid>...</pnrid><email>...</email><sendticket>true</sendticket></request>
// // Response: <response type="AUTOTICKET"><pnrid>...</pnrid><accountno>...</accountno><timelimit>...</timelimit></response>

// type AutoTicketResponse struct {
// 	XMLName   xml.Name `xml:"response"`
// 	PNRID     string   `xml:"pnrid"`
// 	AccountNo string   `xml:"accountno"`
// 	TimeLimit string   `xml:"timelimit"`
// 	ErrorCode string   `xml:"code"`
// 	ErrorMsg  string   `xml:"description"`
// }

// func (c *Client) AutoTicket(pnr, email string, sendTicket bool) (*AutoTicketResponse, error) {
// 	send := "false"
// 	if sendTicket {
// 		send = "true"
// 	}
// 	body := fmt.Sprintf(`<request type="AUTOTICKET"><pnrid>%s</pnrid><email>%s</email><sendticket>%s</sendticket></request>`,
// 		pnr, email, send)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res AutoTicketResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}
// 	return &res, nil
// }

// // ─── CANCEL TICKET ───────────────────────────────────────────────────────────
// //
// // Request:  <request type="CANCEL"><pnrid>...</pnrid></request>
// // Response: <response type="CANCEL"><pnrid>...</pnrid></response>
// // Note: tidak ada field <reason> sesuai dokumentasi Voltras

// type CancelTicketResponse struct {
// 	XMLName   xml.Name `xml:"response"`
// 	PNRID     string   `xml:"pnrid"`
// 	ErrorCode string   `xml:"code"`
// 	ErrorMsg  string   `xml:"description"`
// }

// func (c *Client) CancelTicket(pnr string) (*CancelTicketResponse, error) {
// 	body := fmt.Sprintf(`<request type="CANCEL"><pnrid>%s</pnrid></request>`, pnr)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res CancelTicketResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}
// 	return &res, nil
// }

// // ─── PRINT TICKET ────────────────────────────────────────────────────────────
// //
// // Request:  <request type="PRINT"><pnrid>...</pnrid></request>
// // Response: raw bytes (HTML atau PDF tergantung maskapai)

// func (c *Client) PrintTicket(pnr string) ([]byte, string, error) {
// 	body := fmt.Sprintf(`<request type="PRINT"><pnrid>%s</pnrid></request>`, pnr)

// 	eMessage := generateSignature(body, c.Token)

// 	form := url.Values{}
// 	form.Set("OFFICE_CODE", c.OfficeCode)
// 	form.Set("MESSAGE", body)
// 	form.Set("E_MESSAGE", eMessage)

// 	req, err := http.NewRequest("POST", c.URL, strings.NewReader(form.Encode()))
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	resp, err := c.HTTPClient.Do(req)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer resp.Body.Close()

// 	raw, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode != 200 {
// 		return nil, "", fmt.Errorf("http status %d", resp.StatusCode)
// 	}

// 	contentType := resp.Header.Get("Content-Type")
// 	return raw, contentType, nil
// }

// // ─── PRINT INSURANCE ─────────────────────────────────────────────────────────
// //
// // Request:  <request type="PRINTINSURANCE"><pnrid>...</pnrid></request>
// // Response: binary PDF

// func (c *Client) PrintInsurance(pnr string) ([]byte, error) {
// 	body := fmt.Sprintf(`<request type="PRINTINSURANCE"><pnrid>%s</pnrid></request>`, pnr)

// 	eMessage := generateSignature(body, c.Token)

// 	form := url.Values{}
// 	form.Set("OFFICE_CODE", c.OfficeCode)
// 	form.Set("MESSAGE", body)
// 	form.Set("E_MESSAGE", eMessage)

// 	req, err := http.NewRequest("POST", c.URL, strings.NewReader(form.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	resp, err := c.HTTPClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	raw, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("http status %d", resp.StatusCode)
// 	}
// 	return raw, nil
// }

// // ─── ADVANCE RETRIEVE ────────────────────────────────────────────────────────
// //
// // Request:
// //   <request type="ADVANCERETRIEVE">
// //     <airline>GA</airline>
// //     <fromdate>01-Jan-2025</fromdate>
// //     <todate>01-May-2025</todate>
// //     <status>ALL</status>   <!-- ALL, BOOKED, TICKETED, CANCELED -->
// //   </request>

// type AdvanceRetrievePNR struct {
// 	PNRID       string `xml:"pnrid"`
// 	Airline     string `xml:"airline"`
// 	BookingCode string `xml:"bookingcode"`
// 	Name        string `xml:"name"`
// 	BookingDate string `xml:"bookingdate"`
// 	TotalFare   string `xml:"totalfare"`
// 	NTSA        string `xml:"ntsa"`
// 	Commission  string `xml:"commission"`
// 	Status      string `xml:"status"`
// }

// type AdvanceRetrieveResponse struct {
// 	XMLName   xml.Name             `xml:"response"`
// 	ListPNR   []AdvanceRetrievePNR `xml:"listpnr>pnr"`
// 	ErrorCode string               `xml:"code"`
// 	ErrorMsg  string               `xml:"description"`
// }

// func (c *Client) AdvanceRetrieve(airline, fromDate, toDate, status string) (*AdvanceRetrieveResponse, error) {
// 	body := fmt.Sprintf(`<request type="ADVANCERETRIEVE"><airline>%s</airline><fromdate>%s</fromdate><todate>%s</todate><status>%s</status></request>`,
// 		airline, fromDate, toDate, status)

// 	raw, err := c.post(body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res AdvanceRetrieveResponse
// 	if err := parseXML(raw, &res); err != nil {
// 		return nil, err
// 	}
// 	if res.ErrorCode != "" {
// 		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
// 	}
// 	return &res, nil
// }

package voltras

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"mobile-api/models"
)

// ─── SIGNATURE ───────────────────────────────────────────────────────────────

func generateSignature(message, token string) string {
	h := sha1.New()
	h.Write([]byte(message + token))
	return hex.EncodeToString(h.Sum(nil))
}

// ─── CLIENT ──────────────────────────────────────────────────────────────────

type Client struct {
	URL        string
	OfficeCode string
	Token      string
	HTTPClient *http.Client
}

func NewClient(url, officeCode, token string) *Client {
	return &Client{
		URL:        url,
		OfficeCode: officeCode,
		Token:      token,
		HTTPClient: &http.Client{Timeout: 90 * time.Second},
	}
}

// ─── POST REQUEST ────────────────────────────────────────────────────────────

func (c *Client) post(message string) ([]byte, error) {
	eMessage := generateSignature(message, c.Token)

	form := url.Values{}
	form.Set("OFFICE_CODE", c.OfficeCode)
	form.Set("MESSAGE", message)
	form.Set("E_MESSAGE", eMessage)

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	fmt.Println("=== VOLTRAS REQUEST ===")
	fmt.Println(message)
	fmt.Println("=== VOLTRAS RESPONSE STATUS:", resp.Status, "===")
	fmt.Println(string(raw))
	fmt.Println("=== END ===")

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, string(raw))
	}

	return raw, nil
}

// ─── XML HELPER ──────────────────────────────────────────────────────────────

func parseXML(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("xml parse error: %w", err)
	}
	return nil
}

// ─── OFFICE INFO ─────────────────────────────────────────────────────────────

type OfficeInfoResponse struct {
	OfficeCode string  `xml:"Body>GetOfficeInformationResponse>OfficeCode"`
	OfficeName string  `xml:"Body>GetOfficeInformationResponse>OfficeName"`
	Balance    float64 `xml:"Body>GetOfficeInformationResponse>Balance"`
	Currency   string  `xml:"Body>GetOfficeInformationResponse>Currency"`
}

func (c *Client) GetOfficeInformation() (*OfficeInfoResponse, error) {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
  <Body>
    <GetOfficeInformation>
      <OfficeCode>%s</OfficeCode>
      <Token>%s</Token>
    </GetOfficeInformation>
  </Body>
</Envelope>`, c.OfficeCode, c.Token)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res OfficeInfoResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ─── FLIGHT AVAILABILITY ─────────────────────────────────────────────────────
//
// Response XML structure dari Voltras:
//
// <response type="FLIGHTAVAILABILITY">
//   <trip type="departure">
//     <listitem>
//       <item>
//         <durationHour>...</durationHour>        ← camelCase!
//         <durationMinute>...</durationMinute>    ← camelCase!
//         <throughfare>...</throughfare>
//         <dayexchange>...</dayexchange>
//         <listflight>
//           <flight>
//             <flightno>...</flightno>
//             <airlinecode>...</airlinecode>
//             <fromcity>...</fromcity>
//             <tocity>...</tocity>
//             <departdate>...</departdate>
//             <arrivaldate>...</arrivaldate>
//             <departtime>...</departtime>
//             <arrivetime>...</arrivetime>
//             <durationHour>...</durationHour>    ← camelCase!
//             <durationMinute>...</durationMinute>← camelCase!
//             <dayexchange>...</dayexchange>
//             <meal>...</meal>
//           </flight>
//         </listflight>
//         <listclassgroup>
//           <classgroup>
//             <class>
//               <code>...</code>
//               <availcode>...</availcode>
//               <fare>...</fare>
//               <currency>...</currency>
//               <pairid>...</pairid>
//               <baggagekilos>...</baggagekilos>
//               <cabin>...</cabin>
//               <servicefee>...</servicefee>
//             </class>
//           </classgroup>
//         </listclassgroup>
//       </item>
//     </listitem>
//   </trip>
// </response>

type VolClass struct {
	Code         string  `xml:"code"`
	AvailCode    string  `xml:"availcode"`
	Fare         float64 `xml:"fare"`
	Currency     string  `xml:"currency"`
	PairID       string  `xml:"pairid"`
	BaggageKilos string  `xml:"baggagekilos"`
	Cabin        string  `xml:"cabin"`
	ServiceFee   float64 `xml:"servicefee"`
}

type VolClassGroup struct {
	Classes []VolClass `xml:"class"`
}

type VolFlight struct {
	FlightNo       string `xml:"flightno"`
	AirlineCode    string `xml:"airlinecode"`
	FromCity       string `xml:"fromcity"`
	ToCity         string `xml:"tocity"`
	DepartDate     string `xml:"departdate"`
	ArrivalDate    string `xml:"arrivaldate"`
	DepartTime     string `xml:"departtime"`
	ArriveTime     string `xml:"arrivetime"`
	DurationHour   int    `xml:"durationHour"`   // camelCase sesuai response XML
	DurationMinute int    `xml:"durationMinute"` // camelCase sesuai response XML
	DayExchange    int    `xml:"dayexchange"`
	Meal           bool   `xml:"meal"`
}

type VolFlightItem struct {
	DurationHour   int    `xml:"durationHour"`   // camelCase
	DurationMinute int    `xml:"durationMinute"` // camelCase
	Throughfare    bool   `xml:"throughfare"`
	DayExchange    int    `xml:"dayexchange"`
	ListFlight     struct {
		Flights []VolFlight `xml:"flight"`
	} `xml:"listflight"`
	ListClassGroup []VolClassGroup `xml:"listclassgroup>classgroup"`
}

type VolTrip struct {
	Type     string `xml:"type,attr"` // "departure" atau "return"
	ListItem struct {
		Items []VolFlightItem `xml:"item"`
	} `xml:"listitem"`
}

// Root: <response type="FLIGHTAVAILABILITY">
type SearchFlightResponse struct {
	XMLName  xml.Name  `xml:"response"`
	Trips    []VolTrip `xml:"trip"`
	ErrorMsg string    `xml:"description"`
}

func (c *Client) FlightAvailability(origin, dest, date string, adt, chd, inf int) (*SearchFlightResponse, error) {
	// Parse & validasi tanggal
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (gunakan YYYY-MM-DD): %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	parsed = parsed.In(loc)
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	if parsed.Before(today) {
		return nil, fmt.Errorf("tanggal keberangkatan tidak boleh kurang dari hari ini")
	}

	formattedDate := parsed.Format("02-Jan-2006")

	body := fmt.Sprintf(`<request type="FLIGHTAVAILABILITY">
	<fromcity>%s</fromcity>
	<tocity>%s</tocity>
	<departdate>%s</departdate>
	<listairline>
		<airline>GA</airline>
		<airline>JT</airline>
	</listairline>
	<cabin></cabin>
	<adult>%d</adult>
	<child>%d</child>
	<infant>%d</infant>
	<cheapestclass>true</cheapestclass>
</request>`, origin, dest, formattedDate, adt, chd, inf)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}

	var res SearchFlightResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ─── RETRIEVE FARE (AIRFARE) ─────────────────────────────────────────────────
//
// Request: <request type="AIRFARE"> — bukan SOAP!
// Response:
//   <response type="AIRFARE">
//     <totalfare>
//       <airline code="GA">7752500</airline>
//     </totalfare>
//     <servicefee>
//       <airline code="GA">775250</airline>
//     </servicefee>
//     <insurance>0</insurance>
//   </response>

// AirlineFare mewakili <airline code="GA">7752500</airline>
type AirlineFare struct {
	Code   string  `xml:"code,attr"`
	Amount float64 `xml:",chardata"`
}

// AirfareResponse adalah response dari <request type="AIRFARE">
type AirfareResponse struct {
	XMLName    xml.Name      `xml:"response"`
	TotalFares []AirlineFare `xml:"totalfare>airline"`
	ServiceFees []AirlineFare `xml:"servicefee>airline"`
	Insurance  float64       `xml:"insurance"`
	ErrorCode  string        `xml:"code"`
	ErrorMsg   string        `xml:"description"`
}

// RetrieveFareSegment adalah segment penerbangan untuk AIRFARE request
type RetrieveFareSegment struct {
	FlightNo   string
	ClassCode  string
	FromCity   string
	ToCity     string
	DepartDate string // format: DD-Mon-YYYY (sudah diformat sebelum masuk)
	DepartTime string
	ArriveTime string
}

// RetrieveFareTrip adalah 1 trip (departure / return) untuk AIRFARE request
type RetrieveFareTrip struct {
	TripType    string // "departure" atau "return"
	AirlineCode string
	TripDate    string // format: DD-Mon-YYYY
	Throughfare bool   // true = throughfare, false = sumlocal
	Segments    []RetrieveFareSegment
	ClassCode   string // dipakai saat throughfare = true
}

func (c *Client) RetrieveFare(
	trips []RetrieveFareTrip,
	adultCount, childCount, infantCount int,
	withInsurance bool,
) (*AirfareResponse, error) {

	var sb strings.Builder
	sb.WriteString(`<request type="AIRFARE">`)

	for _, t := range trips {
		sb.WriteString(`<trip type="` + t.TripType + `">`)
		sb.WriteString(`<airlinecode>` + t.AirlineCode + `</airlinecode>`)
		sb.WriteString(`<tripdate>` + t.TripDate + `</tripdate>`)
		sb.WriteString(`<listselecteditem>`)
		sb.WriteString(`<selecteditem>`)

		if t.Throughfare {
			// ── THROUGHFARE ──
			sb.WriteString(`<throughfare>`)
			sb.WriteString(`<listflight>`)
			for _, seg := range t.Segments {
				sb.WriteString(`<segment>`)
				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
				sb.WriteString(`</segment>`)
			}
			sb.WriteString(`</listflight>`)
			sb.WriteString(`<classcode>` + t.ClassCode + `</classcode>`)
			sb.WriteString(`</throughfare>`)
		} else {
			// ── SUM LOCAL (direct atau connecting sum local) ──
			sb.WriteString(`<sumlocal>`)
			for _, seg := range t.Segments {
				sb.WriteString(`<segment>`)
				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
				sb.WriteString(`<classcode>` + seg.ClassCode + `</classcode>`)
				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
				sb.WriteString(`</segment>`)
			}
			sb.WriteString(`</sumlocal>`)
		}

		sb.WriteString(`</selecteditem>`)
		sb.WriteString(`</listselecteditem>`)
		sb.WriteString(`</trip>`)
	}

	ins := "false"
	if withInsurance {
		ins = "true"
	}
	sb.WriteString(`<withinsurance>` + ins + `</withinsurance>`)
	sb.WriteString(fmt.Sprintf(`<adultcount>%d</adultcount>`, adultCount))
	sb.WriteString(fmt.Sprintf(`<childcount>%d</childcount>`, childCount))
	sb.WriteString(fmt.Sprintf(`<infantcount>%d</infantcount>`, infantCount))
	sb.WriteString(`</request>`)

	raw, err := c.post(sb.String())
	if err != nil {
		return nil, err
	}

	var res AirfareResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}

	// Cek error dari Voltras
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}

	return &res, nil
}

// ─── BOOK FLIGHT ─────────────────────────────────────────────────────────────
//
// Request: <request type="BOOK">
// Response:
//   <response type="BOOK">
//     <departure>
//       <pnrid>14773</pnrid>
//       <ntsa>1689467</ntsa>
//       <commission>19633</commission>
//       <totalfare>1709100</totalfare>
//       <servicefee>0</servicefee>
//     </departure>
//     <return>                    ← hanya muncul jika return dengan maskapai BERBEDA
//       <pnrid>...</pnrid>
//       ...
//     </return>
//   </response>

type BookPassenger struct {
	Type           string
	Title          string
	FirstName      string
	LastName       string
	ID             string // KTP
	DOB            string // format: dd-MMM-yyyy (sudah diformat sebelum masuk)
	AdultAssoc     string // pax_id adult untuk INFANT
	Nationality    string
	PassportNo     string
	PassportExp    string // format: dd-MMM-yyyy
	IssuingCountry string
}

// BookTripSegment adalah segment per trip untuk booking
type BookTripSegment struct {
	FlightNo   string
	ClassCode  string // dipakai untuk sumlocal; kosong untuk throughfare
	FromCity   string
	ToCity     string
	DepartDate string // format: DD-Mon-YYYY (sudah diformat)
	DepartTime string
	ArriveTime string
}

// BookTrip adalah 1 trip (departure/return) untuk booking
type BookTrip struct {
	TripType    string // "departure" atau "return"
	AirlineCode string
	TripDate    string            // format: DD-Mon-YYYY
	Throughfare bool              // true = throughfare, false = sumlocal
	ClassCode   string            // classcode untuk throughfare (level trip)
	Segments    []BookTripSegment
}

// BookFlightPNR adalah hasil booking per trip
type BookFlightPNR struct {
	PNRID      string  `xml:"pnrid"`
	NTSA       float64 `xml:"ntsa"`
	Commission float64 `xml:"commission"`
	TotalFare  float64 `xml:"totalfare"`
	ServiceFee float64 `xml:"servicefee"`
}

// BookFlightResponse adalah root response dari Voltras booking
// <response type="BOOK"><departure>...</departure><return>...</return></response>
type BookFlightResponse struct {
	XMLName   xml.Name      `xml:"response"`
	Departure BookFlightPNR `xml:"departure"`
	Return    BookFlightPNR `xml:"return"`    // kosong jika OW atau RT maskapai sama
	ErrorCode string        `xml:"code"`
	ErrorMsg  string        `xml:"description"`
}

func (c *Client) BookFlight(
	trips []BookTrip,
	adultCount, childCount, infantCount int,
	withInsurance bool,
	serviceFee int,
	contact models.PaxContact,
	agentContact models.AgentContact,
	passengers []BookPassenger,
) (*BookFlightResponse, error) {

	var sb strings.Builder

	sb.WriteString(`<request type="BOOK">`)
	sb.WriteString(`<selectionrequest>`)

	// withinsurance wajib ada di dalam selectionrequest (sebelum trip)
	ins := "false"
	if withInsurance {
		ins = "true"
	}
	sb.WriteString(`<withinsurance>` + ins + `</withinsurance>`)
	sb.WriteString(fmt.Sprintf(`<adultcount>%d</adultcount>`, adultCount))
	sb.WriteString(fmt.Sprintf(`<childcount>%d</childcount>`, childCount))
	sb.WriteString(fmt.Sprintf(`<infantcount>%d</infantcount>`, infantCount))

	// ── TRIPS ──
	for _, t := range trips {
		sb.WriteString(`<trip type="` + t.TripType + `">`)
		sb.WriteString(`<airlinecode>` + t.AirlineCode + `</airlinecode>`)
		sb.WriteString(`<tripdate>` + t.TripDate + `</tripdate>`)
		sb.WriteString(`<listselecteditem>`)
		sb.WriteString(`<selecteditem>`)

		if t.Throughfare {
			// ── THROUGHFARE: classcode di level trip, segment tanpa classcode ──
			sb.WriteString(`<throughfare>`)
			for _, seg := range t.Segments {
				sb.WriteString(`<segment>`)
				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
				sb.WriteString(`</segment>`)
			}
			sb.WriteString(`<classcode>` + t.ClassCode + `</classcode>`)
			sb.WriteString(`</throughfare>`)
		} else {
			// ── SUMLOCAL: classcode di tiap segment (direct atau connecting) ──
			sb.WriteString(`<sumlocal>`)
			for _, seg := range t.Segments {
				sb.WriteString(`<segment>`)
				sb.WriteString(`<flightno>` + seg.FlightNo + `</flightno>`)
				sb.WriteString(`<fromcity>` + seg.FromCity + `</fromcity>`)
				sb.WriteString(`<tocity>` + seg.ToCity + `</tocity>`)
				sb.WriteString(`<classcode>` + seg.ClassCode + `</classcode>`)
				sb.WriteString(`<departtime>` + seg.DepartTime + `</departtime>`)
				sb.WriteString(`<arrivetime>` + seg.ArriveTime + `</arrivetime>`)
				sb.WriteString(`<departdate>` + seg.DepartDate + `</departdate>`)
				sb.WriteString(`</segment>`)
			}
			sb.WriteString(`</sumlocal>`)
		}

		sb.WriteString(`</selecteditem>`)
		sb.WriteString(`</listselecteditem>`)
		sb.WriteString(`</trip>`)
	}

	sb.WriteString(`</selectionrequest>`)

	// ── LISTPAX ──
	sb.WriteString(`<listpax>`)
	for _, p := range passengers {
		sb.WriteString(`<pax>`)
		sb.WriteString(`<title>` + p.Title + `</title>`)
		sb.WriteString(`<firstname>` + p.FirstName + `</firstname>`)
		sb.WriteString(`<lastname>` + p.LastName + `</lastname>`)
		sb.WriteString(`<type>` + p.Type + `</type>`)

		// id (KTP) — selalu emit tag meski kosong (sesuai contoh Voltras)
		sb.WriteString(`<id>` + p.ID + `</id>`)

		// dob — selalu emit tag meski kosong
		sb.WriteString(`<dob>` + p.DOB + `</dob>`)

		// adultassoc — untuk INFANT wajib
		sb.WriteString(`<adultassoc>` + p.AdultAssoc + `</adultassoc>`)

		// passport — selalu emit (sesuai contoh Voltras, isi kosong untuk domestik)
		sb.WriteString(`<passport>`)
		sb.WriteString(`<no>` + p.PassportNo + `</no>`)
		sb.WriteString(`<nationality>` + p.Nationality + `</nationality>`)
		sb.WriteString(`<issuingcountry>` + p.IssuingCountry + `</issuingcountry>`)
		sb.WriteString(`<expiry>` + p.PassportExp + `</expiry>`)
		sb.WriteString(`</passport>`)

		sb.WriteString(`</pax>`)
	}
	sb.WriteString(`</listpax>`)

	// ── PAXCONTACT ──
	sb.WriteString(`<paxcontact>`)
	sb.WriteString(`<firstname>` + contact.FirstName + `</firstname>`)
	sb.WriteString(`<lastname>` + contact.LastName + `</lastname>`)
	sb.WriteString(`<phone1>` + contact.Phone1 + `</phone1>`)
	if contact.Phone2 != "" {
		sb.WriteString(`<phone2>` + contact.Phone2 + `</phone2>`)
	}
	if contact.EmailAddress != nil && *contact.EmailAddress != "" {
		sb.WriteString(`<emailaddress>` + *contact.EmailAddress + `</emailaddress>`)
	}
	sb.WriteString(`</paxcontact>`)

	// ── AGENTCONTACT ──
	sb.WriteString(`<agentcontact>`)
	sb.WriteString(`<firstname>` + agentContact.FirstName + `</firstname>`)
	sb.WriteString(`<lastname>` + agentContact.LastName + `</lastname>`)
	sb.WriteString(`<emailaddress>` + agentContact.EmailAddress + `</emailaddress>`)
	if agentContact.Phone != "" {
		sb.WriteString(`<phone>` + agentContact.Phone + `</phone>`)
	}
	sb.WriteString(`</agentcontact>`)

	// ── SERVICE FEE (opsional) ──
	if serviceFee > 0 {
		sb.WriteString(fmt.Sprintf(`<servicefee>%d</servicefee>`, serviceFee))
	}

	sb.WriteString(`</request>`)

	raw, err := c.post(sb.String())
	if err != nil {
		return nil, err
	}

	var res BookFlightResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}

	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}

	return &res, nil
}

// ─── RETRIEVE PNR ────────────────────────────────────────────────────────────
//
// Request:  <request type="RETRIEVE"><pnrid>...</pnrid></request>
// Response: <response type="RETRIEVE"><bookingcode>...</bookingcode>...<listpax>...<listflight>...

type RetrievePaxDetail struct {
	PaxNo      string `xml:"paxno"`
	Title      string `xml:"title"`
	Name       string `xml:"name"`
	TicketNo   string `xml:"ticketno"`
	AdultAssoc string `xml:"adultassoc"`
	DOB        string `xml:"dob"`
}

type RetrieveFlightDetail struct {
	AirlineCode string `xml:"airlinecode"`
	FlightNo    string `xml:"flightno"`
	Departure   string `xml:"departure"`  // kota asal
	Arrival     string `xml:"arrival"`    // kota tujuan
	DepartDate  string `xml:"departdate"`
	DepartTime  string `xml:"departtime"`
	ArrivalTime string `xml:"arrivaltime"`
	ClassCode   string `xml:"classcode"`
	Type        string `xml:"type"` // GO or BACK
}

type RetrievePaxContact struct {
	FirstName string `xml:"firstname"`
	LastName  string `xml:"lastname"`
	Phone1    string `xml:"phone1"`
	Phone2    string `xml:"phone2"`
	Email     string `xml:"email"`
}

type RetrieveNTSA struct {
	Fare      float64 `xml:"fare"`
	Insurance float64 `xml:"insurance"`
}

type RetrievePublish struct {
	Fare       float64 `xml:"fare"`
	Insurance  float64 `xml:"insurance"`
	ServiceFee float64 `xml:"serviceFee"`
}

type RetrieveAutoTicket struct {
	No        string `xml:"no"`
	Amount    string `xml:"amount"`
	TimeLimit string `xml:"timelimit"`
	Status    string `xml:"status"`
}

// RetrievePNRResponse adalah response dari <request type="RETRIEVE">
type RetrievePNRResponse struct {
	XMLName     xml.Name               `xml:"response"`
	BookingCode string                 `xml:"bookingcode"`
	AirlineCode string                 `xml:"airlinecode"`
	Status      string                 `xml:"status"`    // BOOKED, TICKETED, CANCELED
	TimeLimit   string                 `xml:"timelimit"` // tidak null jika BOOKED
	NTSA        RetrieveNTSA           `xml:"ntsa"`
	Publish     RetrievePublish        `xml:"publish"`
	ListPax     []RetrievePaxDetail    `xml:"listpax>pax"`
	PaxContact  RetrievePaxContact     `xml:"paxcontact"`
	ListFlight  []RetrieveFlightDetail `xml:"listflight>flight"`
	ListRemark  []string               `xml:"listremark>remark"`
	AutoTicket  RetrieveAutoTicket     `xml:"autoticket"`
	ErrorCode   string                 `xml:"code"`
	ErrorMsg    string                 `xml:"description"`
}

func (c *Client) RetrievePNR(pnr string) (*RetrievePNRResponse, error) {
	body := fmt.Sprintf(`<request type="RETRIEVE"><pnrid>%s</pnrid></request>`, pnr)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res RetrievePNRResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}
	return &res, nil
}

// ─── TICKETING ───────────────────────────────────────────────────────────────
//
// Request:  <request type="TICKET"><pnrid>...</pnrid></request>
// Response: <response type="TICKET"><pnrid>...</pnrid></response>

type TicketingResponse struct {
	XMLName   xml.Name `xml:"response"`
	PNRID     string   `xml:"pnrid"`
	ErrorCode string   `xml:"code"`
	ErrorMsg  string   `xml:"description"`
}

func (c *Client) TicketingFlight(pnr string) (*TicketingResponse, error) {
	body := fmt.Sprintf(`<request type="TICKET"><pnrid>%s</pnrid></request>`, pnr)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res TicketingResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}
	return &res, nil
}

// ─── AUTO TICKET ─────────────────────────────────────────────────────────────
//
// Request:  <request type="AUTOTICKET"><pnrid>...</pnrid><email>...</email><sendticket>true</sendticket></request>
// Response: <response type="AUTOTICKET"><pnrid>...</pnrid><accountno>...</accountno><timelimit>...</timelimit></response>

type AutoTicketResponse struct {
	XMLName   xml.Name `xml:"response"`
	PNRID     string   `xml:"pnrid"`
	AccountNo string   `xml:"accountno"`
	TimeLimit string   `xml:"timelimit"`
	ErrorCode string   `xml:"code"`
	ErrorMsg  string   `xml:"description"`
}

func (c *Client) AutoTicket(pnr, email string, sendTicket bool) (*AutoTicketResponse, error) {
	send := "false"
	if sendTicket {
		send = "true"
	}
	body := fmt.Sprintf(`<request type="AUTOTICKET"><pnrid>%s</pnrid><email>%s</email><sendticket>%s</sendticket></request>`,
		pnr, email, send)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res AutoTicketResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}
	return &res, nil
}

// ─── CANCEL TICKET ───────────────────────────────────────────────────────────
//
// Request:  <request type="CANCEL"><pnrid>...</pnrid></request>
// Response: <response type="CANCEL"><pnrid>...</pnrid></response>
// Note: tidak ada field <reason> sesuai dokumentasi Voltras

type CancelTicketResponse struct {
	XMLName   xml.Name `xml:"response"`
	PNRID     string   `xml:"pnrid"`
	ErrorCode string   `xml:"code"`
	ErrorMsg  string   `xml:"description"`
}

func (c *Client) CancelTicket(pnr string) (*CancelTicketResponse, error) {
	body := fmt.Sprintf(`<request type="CANCEL"><pnrid>%s</pnrid></request>`, pnr)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res CancelTicketResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}
	return &res, nil
}

// ─── PRINT TICKET ────────────────────────────────────────────────────────────
//
// Request:  <request type="PRINT"><pnrid>...</pnrid></request>
// Response: raw bytes (HTML atau PDF tergantung maskapai)

func (c *Client) PrintTicket(pnr string) ([]byte, string, error) {
	body := fmt.Sprintf(`<request type="PRINT"><pnrid>%s</pnrid></request>`, pnr)

	eMessage := generateSignature(body, c.Token)

	form := url.Values{}
	form.Set("OFFICE_CODE", c.OfficeCode)
	form.Set("MESSAGE", body)
	form.Set("E_MESSAGE", eMessage)

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("http status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	return raw, contentType, nil
}

// ─── PRINT INSURANCE ─────────────────────────────────────────────────────────
//
// Request:  <request type="PRINTINSURANCE"><pnrid>...</pnrid></request>
// Response: binary PDF

func (c *Client) PrintInsurance(pnr string) ([]byte, error) {
	body := fmt.Sprintf(`<request type="PRINTINSURANCE"><pnrid>%s</pnrid></request>`, pnr)

	eMessage := generateSignature(body, c.Token)

	form := url.Values{}
	form.Set("OFFICE_CODE", c.OfficeCode)
	form.Set("MESSAGE", body)
	form.Set("E_MESSAGE", eMessage)

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	return raw, nil
}

// ─── ADVANCE RETRIEVE ────────────────────────────────────────────────────────
//
// Request:
//   <request type="ADVANCERETRIEVE">
//     <airline>GA</airline>
//     <fromdate>01-Jan-2025</fromdate>
//     <todate>01-May-2025</todate>
//     <status>ALL</status>   <!-- ALL, BOOKED, TICKETED, CANCELED -->
//   </request>

type AdvanceRetrievePNR struct {
	PNRID       string `xml:"pnrid"`
	Airline     string `xml:"airline"`
	BookingCode string `xml:"bookingcode"`
	Name        string `xml:"name"`
	BookingDate string `xml:"bookingdate"`
	TotalFare   string `xml:"totalfare"`
	NTSA        string `xml:"ntsa"`
	Commission  string `xml:"commission"`
	Status      string `xml:"status"`
}

type AdvanceRetrieveResponse struct {
	XMLName   xml.Name             `xml:"response"`
	ListPNR   []AdvanceRetrievePNR `xml:"listpnr>pnr"`
	ErrorCode string               `xml:"code"`
	ErrorMsg  string               `xml:"description"`
}

func (c *Client) AdvanceRetrieve(airline, fromDate, toDate, status string) (*AdvanceRetrieveResponse, error) {
	body := fmt.Sprintf(`<request type="ADVANCERETRIEVE"><airline>%s</airline><fromdate>%s</fromdate><todate>%s</todate><status>%s</status></request>`,
		airline, fromDate, toDate, status)

	raw, err := c.post(body)
	if err != nil {
		return nil, err
	}
	var res AdvanceRetrieveResponse
	if err := parseXML(raw, &res); err != nil {
		return nil, err
	}
	if res.ErrorCode != "" {
		return nil, fmt.Errorf("voltras error [%s]: %s", res.ErrorCode, res.ErrorMsg)
	}
	return &res, nil
}