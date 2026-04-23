// package voltras

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"

// 	"mobile-api/models"
// 	"mobile-api/utils"
// )

// var volClient *Client

// func InitClient(url, officeCode, token string) {
// 	volClient = NewClient(url, officeCode, token)
// }

// // GET /api/v1/voltras/office-info
// func GetOfficeInformationHandler(c *gin.Context) {
// 	info, err := volClient.GetOfficeInformation()
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to reach Voltras", err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, utils.SuccessResponse("Office information retrieved", gin.H{
// 		"office_code": info.OfficeCode,
// 		"office_name": info.OfficeName,
// 		"balance":     info.Balance,
// 		"currency":    info.Currency,
// 	}))
// }

// // POST /api/v1/voltras/flight/search
// func FlightAvailabilityHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var req models.FlightAvailabilityRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		// Validasi field wajib
// 		if req.FromCity == "" || req.ToCity == "" || req.DepartDate == "" {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 				"Missing required fields",
// 				"fromcity, tocity, departdate wajib diisi",
// 			))
// 			return
// 		}

// 		// Default passenger
// 		adult, _ := strconv.Atoi(req.Adult)
// 		if adult == 0 {
// 			adult = 1
// 		}
// 		child, _ := strconv.Atoi(req.Child)
// 		infant, _ := strconv.Atoi(req.Infant)

// 		// Call Voltras
// 		result, err := volClient.FlightAvailability(
// 			req.FromCity,
// 			req.ToCity,
// 			req.DepartDate,
// 			adult,
// 			child,
// 			infant,
// 		)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Search failed", err.Error()))
// 			return
// 		}

// 		if result.ErrorMsg != "" {
// 			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse("Voltras error", result.ErrorMsg))
// 			return
// 		}

// 		// Kumpulkan semua flight dari trip departure
// 		// Setiap item = 1 opsi itinerary (bisa direct atau connecting)
// 		// Setiap item punya listflight (list segment) dan listclassgroup (list kelas)
// 		type FlightOption struct {
// 			DurationHour   int           `json:"duration_hour"`
// 			DurationMinute int           `json:"duration_minute"`
// 			Throughfare    bool          `json:"throughfare"`
// 			DayExchange    int           `json:"day_exchange"`
// 			Segments       []VolFlight   `json:"segments"`
// 			Classes        []VolClass    `json:"classes"`
// 		}

// 		var options []FlightOption

// 		for _, trip := range result.Trips {
// 			if trip.Type != "departure" {
// 				continue
// 			}
// 			for _, item := range trip.ListItem.Items {
// 				// Kumpulkan semua class dari semua classgroup
// 				var allClasses []VolClass
// 				for _, cg := range item.ListClassGroup {
// 					allClasses = append(allClasses, cg.Classes...)
// 				}

// 				opt := FlightOption{
// 					DurationHour:   item.DurationHour,
// 					DurationMinute: item.DurationMinute,
// 					Throughfare:    item.Throughfare,
// 					DayExchange:    item.DayExchange,
// 					Segments:       item.ListFlight.Flights,
// 					Classes:        allClasses,
// 				}
// 				options = append(options, opt)
// 			}
// 		}

// 		if options == nil {
// 			options = []FlightOption{}
// 		}

// 		c.JSON(http.StatusOK, utils.SuccessResponse("Flights retrieved", gin.H{
// 			"flights":   options,
// 			"total":     len(options),
// 			"trip_type": "OW",
// 		}))
// 	}
// }

// // POST /api/v1/voltras/flight/fare/retrieve
// func RetrieveFareHandler(c *gin.Context) {
// 	var req models.RetrieveFareRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 		return
// 	}

// 	// Validasi minimal ada 1 trip dengan 1 segment
// 	if len(req.Trips) == 0 {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("trips wajib diisi minimal 1", "invalid request"))
// 		return
// 	}

// 	// Convert request trips → RetrieveFareTrip
// 	fareTrips := make([]RetrieveFareTrip, 0, len(req.Trips))
// 	for i, t := range req.Trips {
// 		if len(t.Segments) == 0 {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 				fmt.Sprintf("trips[%d]: segments wajib diisi minimal 1", i),
// 				"invalid request",
// 			))
// 			return
// 		}

// 		// Parse & format trip date ke DD-Mon-YYYY
// 		tripDateParsed, err := time.Parse("2006-01-02", t.TripDate)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 				fmt.Sprintf("trips[%d]: trip_date format tidak valid (gunakan YYYY-MM-DD)", i),
// 				err.Error(),
// 			))
// 			return
// 		}
// 		formattedTripDate := tripDateParsed.Format("02-Jan-2006")

// 		// Convert segments
// 		fareSegs := make([]RetrieveFareSegment, 0, len(t.Segments))
// 		for j, seg := range t.Segments {
// 			// Parse & format depart date ke DD-Mon-YYYY
// 			departDateParsed, err := time.Parse("2006-01-02", seg.DepartDate)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 					fmt.Sprintf("trips[%d].segments[%d]: depart_date format tidak valid", i, j),
// 					err.Error(),
// 				))
// 				return
// 			}
// 			fareSegs = append(fareSegs, RetrieveFareSegment{
// 				FlightNo:   seg.FlightNo,
// 				ClassCode:  seg.ClassCode,
// 				FromCity:   seg.FromCity,
// 				ToCity:     seg.ToCity,
// 				DepartDate: departDateParsed.Format("02-Jan-2006"),
// 				DepartTime: seg.DepartTime,
// 				ArriveTime: seg.ArriveTime,
// 			})
// 		}

// 		fareTrips = append(fareTrips, RetrieveFareTrip{
// 			TripType:    t.TripType,
// 			AirlineCode: t.AirlineCode,
// 			TripDate:    formattedTripDate,
// 			Throughfare: t.Throughfare,
// 			ClassCode:   t.ClassCode,
// 			Segments:    fareSegs,
// 		})
// 	}

// 	// Default passenger count
// 	adult := req.AdultCount
// 	if adult == 0 {
// 		adult = 1
// 	}

// 	result, err := volClient.RetrieveFare(fareTrips, adult, req.ChildCount, req.InfantCount, req.WithInsurance)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Retrieve fare failed", err.Error()))
// 		return
// 	}

// 	// Build response yang lebih readable
// 	type FarePerAirline struct {
// 		AirlineCode string  `json:"airline_code"`
// 		TotalFare   float64 `json:"total_fare"`
// 		ServiceFee  float64 `json:"service_fee"`
// 	}

// 	// Map airline code → fare+servicefee
// 	fareMap := map[string]*FarePerAirline{}
// 	for _, tf := range result.TotalFares {
// 		if _, ok := fareMap[tf.Code]; !ok {
// 			fareMap[tf.Code] = &FarePerAirline{AirlineCode: tf.Code}
// 		}
// 		fareMap[tf.Code].TotalFare = tf.Amount
// 	}
// 	for _, sf := range result.ServiceFees {
// 		if _, ok := fareMap[sf.Code]; !ok {
// 			fareMap[sf.Code] = &FarePerAirline{AirlineCode: sf.Code}
// 		}
// 		fareMap[sf.Code].ServiceFee = sf.Amount
// 	}

// 	fares := make([]FarePerAirline, 0, len(fareMap))
// 	var grandTotal float64
// 	for _, f := range fareMap {
// 		fares = append(fares, *f)
// 		grandTotal += f.TotalFare
// 	}

// 	c.JSON(http.StatusOK, utils.SuccessResponse("Fare retrieved", gin.H{
// 		"fares":       fares,
// 		"insurance":   result.Insurance,
// 		"grand_total": grandTotal + result.Insurance,
// 	}))
// }

// // POST /api/v1/voltras/flight/book
// func BookFlightHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var req models.BookFlightRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		// Default passenger count
// 		if req.AdultCount == 0 {
// 			req.AdultCount = 1
// 		}

// 		// Validasi contact
// 		if req.PaxContact.Phone1 == "" {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("paxcontact.phone1 wajib diisi", "invalid request"))
// 			return
// 		}
// 		if req.AgentContact.EmailAddress == "" {
// 			req.AgentContact = models.AgentContact{
// 				FirstName:    "API",
// 				LastName:     "SYSTEM",
// 				EmailAddress: "api@system.com",
// 			}
// 		}

// 		// ── BUILD TRIPS ──
// 		bookTrips := make([]BookTrip, 0, len(req.Trips))
// 		for i, t := range req.Trips {
// 			if len(t.Segments) == 0 {
// 				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 					fmt.Sprintf("trips[%d]: segments wajib diisi minimal 1", i), "invalid request",
// 				))
// 				return
// 			}

// 			tripDateParsed, err := time.Parse("2006-01-02", t.TripDate)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 					fmt.Sprintf("trips[%d]: trip_date tidak valid (YYYY-MM-DD)", i), err.Error(),
// 				))
// 				return
// 			}
// 			formattedTripDate := tripDateParsed.Format("02-Jan-2006")

// 			bookSegs := make([]BookTripSegment, 0, len(t.Segments))
// 			for j, seg := range t.Segments {
// 				departDateParsed, err := time.Parse("2006-01-02", seg.DepartDate)
// 				if err != nil {
// 					c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 						fmt.Sprintf("trips[%d].segments[%d]: depart_date tidak valid", i, j), err.Error(),
// 					))
// 					return
// 				}
// 				bookSegs = append(bookSegs, BookTripSegment{
// 					FlightNo:   seg.FlightNo,
// 					ClassCode:  seg.ClassCode,
// 					FromCity:   seg.FromCity,
// 					ToCity:     seg.ToCity,
// 					DepartDate: departDateParsed.Format("02-Jan-2006"),
// 					DepartTime: seg.DepartTime,
// 					ArriveTime: seg.ArriveTime,
// 				})
// 			}

// 			bookTrips = append(bookTrips, BookTrip{
// 				TripType:    t.TripType,
// 				AirlineCode: t.AirlineCode,
// 				TripDate:    formattedTripDate,
// 				Throughfare: t.Throughfare,
// 				ClassCode:   t.ClassCode,
// 				Segments:    bookSegs,
// 			})
// 		}

// 		// ── BUILD PASSENGERS ──
// 		passengers := make([]BookPassenger, 0, len(req.ListPax))
// 		for i, p := range req.ListPax {
// 			if p.FirstName == "" {
// 				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 					fmt.Sprintf("listpax[%d]: first_name wajib diisi", i), "invalid passenger",
// 				))
// 				return
// 			}
// 			if p.LastName == "" {
// 				p.LastName = "/"
// 			}

// 			// Format DOB: YYYY-MM-DD → dd-MMM-yyyy (contoh: 10-AUG-1994)
// 			var dob string
// 			if p.DOB != nil && *p.DOB != "" {
// 				t, err := time.Parse("2006-01-02", *p.DOB)
// 				if err == nil {
// 					dob = strings.ToUpper(t.Format("02-Jan-2006"))
// 				} else {
// 					dob = *p.DOB
// 				}
// 			}

// 			// Adult assoc (untuk INFANT)
// 			var adultAssoc string
// 			if p.AdultAssoc != nil {
// 				adultAssoc = fmt.Sprintf("%d", *p.AdultAssoc)
// 			}

// 			// Passport fields
// 			var passportNo, nationality, issuingCountry, passportExp string
// 			if p.Passport != nil {
// 				passportNo = p.Passport.No
// 				nationality = p.Passport.Nationality
// 				issuingCountry = p.Passport.IssuingCountry
// 				if p.Passport.Expiry != "" {
// 					t, err := time.Parse("2006-01-02", p.Passport.Expiry)
// 					if err == nil {
// 						passportExp = strings.ToUpper(t.Format("02-Jan-2006"))
// 					} else {
// 						passportExp = p.Passport.Expiry
// 					}
// 				}
// 			}

// 			passengers = append(passengers, BookPassenger{
// 				Type:           p.Type,
// 				Title:          p.Title,
// 				FirstName:      p.FirstName,
// 				LastName:       p.LastName,
// 				DOB:            dob,
// 				AdultAssoc:     adultAssoc,
// 				Nationality:    nationality,
// 				PassportNo:     passportNo,
// 				IssuingCountry: issuingCountry,
// 				PassportExp:    passportExp,
// 			})
// 		}

// 		// ── HIT VOLTRAS ──
// 		result, err := volClient.BookFlight(
// 			bookTrips,
// 			req.AdultCount,
// 			req.ChildCount,
// 			req.InfantCount,
// 			req.WithInsurance,
// 			req.ServiceFee,
// 			req.PaxContact,
// 			req.AgentContact,
// 			passengers,
// 		)
// 		if err != nil {
// 			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Book flight failed", err.Error()))
// 			return
// 		}

// 		// ── SIMPAN CONTACT ──
// 		contact := models.Contact{
// 			FirstName: req.PaxContact.FirstName,
// 			LastName:  req.PaxContact.LastName,
// 			Phone:     req.PaxContact.Phone1,
// 		}
// 		if req.PaxContact.EmailAddress != nil {
// 			contact.Email = *req.PaxContact.EmailAddress
// 		}
// 		db.Where("phone = ?", contact.Phone).FirstOrCreate(&contact)

// 		// ── SIMPAN TRANSACTION ──
// 		tripType := "oneway"
// 		if len(req.Trips) > 1 {
// 			tripType = "roundtrip"
// 		}
// 		totalAmount := result.Departure.TotalFare + result.Return.TotalFare

// 		flightTx := models.FlightTransaction{
// 			ContactID:      contact.ID,
// 			VMCode:         "CGK001",
// 			Amount:         totalAmount,
// 			TripType:       tripType,
// 			Status:         "pending",
// 			SourcePlatform: "mobile",
// 		}
// 		db.Create(&flightTx)

// 		// ── SIMPAN BOOKING DEPARTURE ──
// 		depBooking := models.FlightBooking{
// 			TransactionID: flightTx.ID,
// 			PNRID:         result.Departure.PNRID,
// 			NTSA:          result.Departure.NTSA,
// 			Commission:    result.Departure.Commission,
// 			PublishFare:   result.Departure.TotalFare,
// 			Status:        "booked",
// 		}
// 		db.Create(&depBooking)

// 		// ── SIMPAN BOOKING RETURN (jika ada PNR return — beda maskapai) ──
// 		var retBooking *models.FlightBooking
// 		if result.Return.PNRID != "" {
// 			rb := models.FlightBooking{
// 				TransactionID: flightTx.ID,
// 				PNRID:         result.Return.PNRID,
// 				NTSA:          result.Return.NTSA,
// 				Commission:    result.Return.Commission,
// 				PublishFare:   result.Return.TotalFare,
// 				Status:        "booked",
// 			}
// 			db.Create(&rb)
// 			retBooking = &rb
// 		}

// 		// ── SIMPAN PASSENGERS ke booking departure ──
// 		for i, p := range req.ListPax {
// 			paxID := i + 1

// 			var dob *time.Time
// 			if p.DOB != nil && *p.DOB != "" {
// 				t, err := time.Parse("2006-01-02", *p.DOB)
// 				if err == nil {
// 					dob = &t
// 				}
// 			}

// 			var passportExp *time.Time
// 			var passportNo, nationality, issuing string
// 			if p.Passport != nil {
// 				passportNo = p.Passport.No
// 				nationality = p.Passport.Nationality
// 				issuing = p.Passport.IssuingCountry
// 				if p.Passport.Expiry != "" {
// 					t, err := time.Parse("2006-01-02", p.Passport.Expiry)
// 					if err == nil {
// 						passportExp = &t
// 					}
// 				}
// 			}

// 			lastName := p.LastName
// 			if lastName == "" {
// 				lastName = "/"
// 			}

// 			db.Create(&models.FlightPassenger{
// 				BookingID:              depBooking.ID,
// 				Title:                  p.Title,
// 				FullName:               p.FirstName + " " + lastName,
// 				Type:                   p.Type,
// 				PaxID:                  paxID,
// 				DOB:                    dob,
// 				AdultAssoc:             p.AdultAssoc,
// 				PassportNo:             passportNo,
// 				PassportNationality:    nationality,
// 				PassportIssuingCountry: issuing,
// 				PassportExpiry:         passportExp,
// 			})
// 		}

// 		// ── RESPONSE ──
// 		resp := gin.H{
// 			"transaction_id": flightTx.ID,
// 			"trip_type":      tripType,
// 			"departure": gin.H{
// 				"booking_id": depBooking.ID,
// 				"pnr":        depBooking.PNRID,
// 				"ntsa":       depBooking.NTSA,
// 				"commission": depBooking.Commission,
// 				"total_fare": depBooking.PublishFare,
// 				"status":     depBooking.Status,
// 			},
// 		}
// 		if retBooking != nil {
// 			resp["return"] = gin.H{
// 				"booking_id": retBooking.ID,
// 				"pnr":        retBooking.PNRID,
// 				"ntsa":       retBooking.NTSA,
// 				"commission": retBooking.Commission,
// 				"total_fare": retBooking.PublishFare,
// 				"status":     retBooking.Status,
// 			}
// 		}

// 		c.JSON(http.StatusCreated, utils.SuccessResponse("Flight booked", resp))
// 	}
// }


// // POST /api/v1/voltras/flight/pnr/retrieve
// func RetrievePNRHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req models.RetrievePNRRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		result, err := volClient.RetrievePNR(req.PNR)
// 		if err != nil {
// 			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Retrieve PNR failed", err.Error()))
// 			return
// 		}

// 		// Sync status + booking_code ke DB
// 		db.Model(&models.FlightBooking{}).
// 			Where("pnrid = ?", req.PNR).
// 			Updates(map[string]interface{}{
// 				"status":       strings.ToLower(result.Status),
// 				"booking_code": result.BookingCode,
// 			})

// 		// Update ticket_number per pax jika sudah TICKETED
// 		if strings.EqualFold(result.Status, "TICKETED") {
// 			for _, pax := range result.ListPax {
// 				if pax.TicketNo != "" {
// 					db.Model(&models.FlightPassenger{}).
// 						Joins("JOIN t_flight_bookings ON t_flight_bookings.id = t_flight_passengers.booking_id").
// 						Where("t_flight_bookings.pnrid = ? AND t_flight_passengers.pax_id = ?", req.PNR, pax.PaxNo).
// 						Update("ticket_number", pax.TicketNo)
// 				}
// 			}
// 		}

// 		c.JSON(http.StatusOK, utils.SuccessResponse("PNR retrieved", result))
// 	}
// }

// // POST /api/v1/voltras/flight/ticket
// func TicketingFlightHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req models.TicketingRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		var booking models.FlightBooking
// 		if err := db.First(&booking, req.BookingID).Error; err != nil {
// 			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
// 			return
// 		}

// 		result, err := volClient.TicketingFlight(booking.PNRID)
// 		if err != nil {
// 			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Ticketing failed", err.Error()))
// 			return
// 		}

// 		db.Model(&booking).Update("status", "ticketed")

// 		c.JSON(http.StatusOK, utils.SuccessResponse("Ticketing successful", gin.H{
// 			"booking_id": booking.ID,
// 			"pnr":        result.PNRID,
// 			"status":     "ticketed",
// 		}))
// 	}
// }

// // POST /api/v1/voltras/flight/autoticket
// func AutoTicketHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req models.AutoTicketRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		var booking models.FlightBooking
// 		if err := db.First(&booking, req.BookingID).Error; err != nil {
// 			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
// 			return
// 		}

// 		if booking.Status != "booked" {
// 			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse(
// 				"Booking status tidak valid untuk auto ticket",
// 				fmt.Sprintf("status saat ini: %s", booking.Status),
// 			))
// 			return
// 		}

// 		result, err := volClient.AutoTicket(booking.PNRID, req.Email, req.SendTicket)
// 		if err != nil {
// 			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Auto ticket failed", err.Error()))
// 			return
// 		}

// 		db.Create(&models.FlightRemark{
// 			BookingID: booking.ID,
// 			Remark:    fmt.Sprintf("AUTOTICKET VA: %s | TIMELIMIT: %s", result.AccountNo, result.TimeLimit),
// 		})

// 		c.JSON(http.StatusOK, utils.SuccessResponse("Auto ticket created", gin.H{
// 			"booking_id": booking.ID,
// 			"pnr":        result.PNRID,
// 			"account_no": result.AccountNo,
// 			"time_limit": result.TimeLimit,
// 		}))
// 	}
// }

// // POST /api/v1/voltras/flight/ticket/cancel
// func CancelTicketHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req models.CancelTicketRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 			return
// 		}

// 		var booking models.FlightBooking
// 		if err := db.First(&booking, req.BookingID).Error; err != nil {
// 			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
// 			return
// 		}

// 		if booking.Status == "cancelled" {
// 			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse("Booking sudah dibatalkan", ""))
// 			return
// 		}

// 		result, err := volClient.CancelTicket(booking.PNRID)
// 		if err != nil {
// 			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Cancel failed", err.Error()))
// 			return
// 		}

// 		db.Model(&booking).Update("status", "cancelled")

// 		remark := "CANCELLED"
// 		if req.Reason != "" {
// 			remark = "CANCELLED: " + req.Reason
// 		}
// 		db.Create(&models.FlightRemark{
// 			BookingID: booking.ID,
// 			Remark:    remark,
// 		})

// 		c.JSON(http.StatusOK, utils.SuccessResponse("Ticket cancelled successfully", gin.H{
// 			"booking_id": booking.ID,
// 			"pnr":        result.PNRID,
// 			"status":     "cancelled",
// 		}))
// 	}
// }

// // GET /api/v1/voltras/flight/print?pnr=XXX
// func PrintTicketHandler(c *gin.Context) {
// 	pnr := c.Query("pnr")
// 	if pnr == "" {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("query param 'pnr' wajib diisi", ""))
// 		return
// 	}

// 	raw, contentType, err := volClient.PrintTicket(pnr)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Print ticket failed", err.Error()))
// 		return
// 	}

// 	filename := "ticket_" + pnr + ".html"
// 	if strings.Contains(contentType, "pdf") {
// 		filename = "ticket_" + pnr + ".pdf"
// 	}

// 	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
// 	c.Data(http.StatusOK, contentType, raw)
// }

// // GET /api/v1/voltras/flight/print-insurance?pnr=XXX
// func PrintInsuranceHandler(c *gin.Context) {
// 	pnr := c.Query("pnr")
// 	if pnr == "" {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("query param 'pnr' wajib diisi", ""))
// 		return
// 	}

// 	raw, err := volClient.PrintInsurance(pnr)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Print insurance failed", err.Error()))
// 		return
// 	}

// 	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="insurance_%s.pdf"`, pnr))
// 	c.Data(http.StatusOK, "application/pdf", raw)
// }

// // POST /api/v1/voltras/flight/advance-retrieve
// func AdvanceRetrieveHandler(c *gin.Context) {
// 	var req models.AdvanceRetrieveRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
// 		return
// 	}

// 	if req.Status == "" {
// 		req.Status = "ALL"
// 	}

// 	validStatus := map[string]bool{"ALL": true, "BOOKED": true, "TICKETED": true, "CANCELED": true}
// 	if !validStatus[strings.ToUpper(req.Status)] {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse(
// 			"status tidak valid, gunakan: ALL, BOOKED, TICKETED, CANCELED", "",
// 		))
// 		return
// 	}

// 	fromParsed, err := time.Parse("2006-01-02", req.FromDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("from_date tidak valid (YYYY-MM-DD)", err.Error()))
// 		return
// 	}
// 	toParsed, err := time.Parse("2006-01-02", req.ToDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("to_date tidak valid (YYYY-MM-DD)", err.Error()))
// 		return
// 	}

// 	result, err := volClient.AdvanceRetrieve(
// 		req.Airline,
// 		fromParsed.Format("02-Jan-2006"),
// 		toParsed.Format("02-Jan-2006"),
// 		strings.ToUpper(req.Status),
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Advance retrieve failed", err.Error()))
// 		return
// 	}

// 	if result.ListPNR == nil {
// 		result.ListPNR = []AdvanceRetrievePNR{}
// 	}

// 	c.JSON(http.StatusOK, utils.SuccessResponse("Advance retrieve success", gin.H{
// 		"total": len(result.ListPNR),
// 		"list":  result.ListPNR,
// 	}))
// }
package voltras

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mobile-api/models"
	"mobile-api/utils"
)

var volClient *Client

func InitClient(url, officeCode, token string) {
	volClient = NewClient(url, officeCode, token)
}

// formatVolDate: YYYY-MM-DD → dd-Mon-yyyy Title Case (contoh: 23-Apr-2026)
// Ini format yang diterima Voltras untuk departdate, tripdate, dob, passport expiry
func formatVolDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}
	// Go time.Format("02-Jan-2006") menghasilkan Title Case bulan: Jan, Feb, Mar, Apr, dst
	return t.Format("02-Jan-2006"), nil
}

// GET /api/v1/voltras/office-info
func GetOfficeInformationHandler(c *gin.Context) {
	info, err := volClient.GetOfficeInformation()
	if err != nil {
		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to reach Voltras", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Office information retrieved", gin.H{
		"office_code": info.OfficeCode,
		"office_name": info.OfficeName,
		"balance":     info.Balance,
		"currency":    info.Currency,
	}))
}

// POST /api/v1/voltras/flight/search
func FlightAvailabilityHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.FlightAvailabilityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		if req.FromCity == "" || req.ToCity == "" || req.DepartDate == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
				"Missing required fields",
				"fromcity, tocity, departdate wajib diisi",
			))
			return
		}

		adult, _ := strconv.Atoi(req.Adult)
		if adult == 0 {
			adult = 1
		}
		child, _ := strconv.Atoi(req.Child)
		infant, _ := strconv.Atoi(req.Infant)

		result, err := volClient.FlightAvailability(req.FromCity, req.ToCity, req.DepartDate, adult, child, infant)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Search failed", err.Error()))
			return
		}

		if result.ErrorMsg != "" {
			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse("Voltras error", result.ErrorMsg))
			return
		}

		type FlightOption struct {
			DurationHour   int         `json:"duration_hour"`
			DurationMinute int         `json:"duration_minute"`
			Throughfare    bool        `json:"throughfare"`
			DayExchange    int         `json:"day_exchange"`
			Segments       []VolFlight `json:"segments"`
			Classes        []VolClass  `json:"classes"`
		}

		var options []FlightOption
		for _, trip := range result.Trips {
			if trip.Type != "departure" {
				continue
			}
			for _, item := range trip.ListItem.Items {
				var allClasses []VolClass
				for _, cg := range item.ListClassGroup {
					allClasses = append(allClasses, cg.Classes...)
				}
				options = append(options, FlightOption{
					DurationHour:   item.DurationHour,
					DurationMinute: item.DurationMinute,
					Throughfare:    item.Throughfare,
					DayExchange:    item.DayExchange,
					Segments:       item.ListFlight.Flights,
					Classes:        allClasses,
				})
			}
		}

		if options == nil {
			options = []FlightOption{}
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Flights retrieved", gin.H{
			"flights":   options,
			"total":     len(options),
			"trip_type": "OW",
		}))
	}
}

// POST /api/v1/voltras/flight/fare/retrieve
func RetrieveFareHandler(c *gin.Context) {
	var req models.RetrieveFareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	if len(req.Trips) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("trips wajib diisi minimal 1", "invalid request"))
		return
	}

	fareTrips := make([]RetrieveFareTrip, 0, len(req.Trips))
	for i, t := range req.Trips {
		if len(t.Segments) == 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
				fmt.Sprintf("trips[%d]: segments wajib diisi minimal 1", i), "invalid request",
			))
			return
		}

		tripDate, err := formatVolDate(t.TripDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
				fmt.Sprintf("trips[%d]: trip_date format tidak valid (gunakan YYYY-MM-DD)", i), err.Error(),
			))
			return
		}

		fareSegs := make([]RetrieveFareSegment, 0, len(t.Segments))
		for j, seg := range t.Segments {
			departDate, err := formatVolDate(seg.DepartDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
					fmt.Sprintf("trips[%d].segments[%d]: depart_date format tidak valid", i, j), err.Error(),
				))
				return
			}
			fareSegs = append(fareSegs, RetrieveFareSegment{
				FlightNo:   seg.FlightNo,
				ClassCode:  seg.ClassCode,
				FromCity:   seg.FromCity,
				ToCity:     seg.ToCity,
				DepartDate: departDate,
				DepartTime: seg.DepartTime,
				ArriveTime: seg.ArriveTime,
			})
		}

		fareTrips = append(fareTrips, RetrieveFareTrip{
			TripType:    t.TripType,
			AirlineCode: t.AirlineCode,
			TripDate:    tripDate,
			Throughfare: t.Throughfare,
			ClassCode:   t.ClassCode,
			Segments:    fareSegs,
		})
	}

	adult := req.AdultCount
	if adult == 0 {
		adult = 1
	}

	result, err := volClient.RetrieveFare(fareTrips, adult, req.ChildCount, req.InfantCount, req.WithInsurance)
	if err != nil {
		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Retrieve fare failed", err.Error()))
		return
	}

	type FarePerAirline struct {
		AirlineCode string  `json:"airline_code"`
		TotalFare   float64 `json:"total_fare"`
		ServiceFee  float64 `json:"service_fee"`
	}

	fareMap := map[string]*FarePerAirline{}
	for _, tf := range result.TotalFares {
		if _, ok := fareMap[tf.Code]; !ok {
			fareMap[tf.Code] = &FarePerAirline{AirlineCode: tf.Code}
		}
		fareMap[tf.Code].TotalFare = tf.Amount
	}
	for _, sf := range result.ServiceFees {
		if _, ok := fareMap[sf.Code]; !ok {
			fareMap[sf.Code] = &FarePerAirline{AirlineCode: sf.Code}
		}
		fareMap[sf.Code].ServiceFee = sf.Amount
	}

	fares := make([]FarePerAirline, 0, len(fareMap))
	var grandTotal float64
	for _, f := range fareMap {
		fares = append(fares, *f)
		grandTotal += f.TotalFare
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Fare retrieved", gin.H{
		"fares":       fares,
		"insurance":   result.Insurance,
		"grand_total": grandTotal + result.Insurance,
	}))
}

// POST /api/v1/voltras/flight/book
func BookFlightHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.BookFlightRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		if req.AdultCount == 0 {
			req.AdultCount = 1
		}
		if req.PaxContact.Phone1 == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("paxcontact.phone1 wajib diisi", "invalid request"))
			return
		}
		if req.AgentContact.EmailAddress == "" {
			req.AgentContact = models.AgentContact{
				FirstName:    "API",
				LastName:     "SYSTEM",
				EmailAddress: "api@system.com",
			}
		}

		// ── BUILD TRIPS ──
		bookTrips := make([]BookTrip, 0, len(req.Trips))
		for i, t := range req.Trips {
			if len(t.Segments) == 0 {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
					fmt.Sprintf("trips[%d]: segments wajib diisi minimal 1", i), "invalid request",
				))
				return
			}

			// FIX: Title Case dd-Mon-yyyy bukan uppercase
			tripDate, err := formatVolDate(t.TripDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
					fmt.Sprintf("trips[%d]: trip_date tidak valid (YYYY-MM-DD)", i), err.Error(),
				))
				return
			}

			bookSegs := make([]BookTripSegment, 0, len(t.Segments))
			for j, seg := range t.Segments {
				// FIX: Title Case dd-Mon-yyyy bukan uppercase
				departDate, err := formatVolDate(seg.DepartDate)
				if err != nil {
					c.JSON(http.StatusBadRequest, utils.ErrorResponse(
						fmt.Sprintf("trips[%d].segments[%d]: depart_date tidak valid", i, j), err.Error(),
					))
					return
				}
				bookSegs = append(bookSegs, BookTripSegment{
					FlightNo:   seg.FlightNo,
					ClassCode:  seg.ClassCode,
					FromCity:   seg.FromCity,
					ToCity:     seg.ToCity,
					DepartDate: departDate, // 23-Apr-2026 bukan 23-APR-2026
					DepartTime: seg.DepartTime,
					ArriveTime: seg.ArriveTime,
				})
			}

			bookTrips = append(bookTrips, BookTrip{
				TripType:    t.TripType,
				AirlineCode: t.AirlineCode,
				TripDate:    tripDate, // 23-Apr-2026 bukan 23-APR-2026
				Throughfare: t.Throughfare,
				ClassCode:   t.ClassCode,
				Segments:    bookSegs,
			})
		}

		// ── BUILD PASSENGERS ──
		passengers := make([]BookPassenger, 0, len(req.ListPax))
		for i, p := range req.ListPax {
			if p.FirstName == "" {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(
					fmt.Sprintf("listpax[%d]: first_name wajib diisi", i), "invalid passenger",
				))
				return
			}
			if p.LastName == "" {
				p.LastName = "/"
			}

			// DOB: FIX Title Case dd-Mon-yyyy
			var dob string
			if p.DOB != nil && *p.DOB != "" {
				formatted, err := formatVolDate(*p.DOB)
				if err == nil {
					dob = formatted // 10-Aug-1994
				} else {
					dob = *p.DOB
				}
			}

			var adultAssoc string
			if p.AdultAssoc != nil {
				adultAssoc = fmt.Sprintf("%d", *p.AdultAssoc)
			}

			// Passport wajib ada (isi kosong untuk domestik)
			var passportNo, nationality, issuingCountry, passportExp string
			if p.Passport != nil {
				passportNo = p.Passport.No
				nationality = p.Passport.Nationality
				issuingCountry = p.Passport.IssuingCountry
				if p.Passport.Expiry != "" {
					// FIX Title Case dd-Mon-yyyy
					formatted, err := formatVolDate(p.Passport.Expiry)
					if err == nil {
						passportExp = formatted
					} else {
						passportExp = p.Passport.Expiry
					}
				}
			}

			passengers = append(passengers, BookPassenger{
				Type:           p.Type,
				Title:          p.Title,
				FirstName:      p.FirstName,
				LastName:       p.LastName,
				DOB:            dob,
				AdultAssoc:     adultAssoc,
				Nationality:    nationality,
				PassportNo:     passportNo,
				IssuingCountry: issuingCountry,
				PassportExp:    passportExp,
			})
		}

		// ── HIT VOLTRAS ──
		result, err := volClient.BookFlight(
			bookTrips,
			req.AdultCount,
			req.ChildCount,
			req.InfantCount,
			req.WithInsurance,
			req.ServiceFee,
			req.PaxContact,
			req.AgentContact,
			passengers,
		)

		if err != nil {
			log.Printf("[BookFlight] Voltras error: %v", err)

		// ─────────────────────────────────────────
		// 🔥 TAMBAHAN: SIMPAN DATA WALAU GAGAL
		// ─────────────────────────────────────────

		// SIMPAN CONTACT
		contact := models.Contact{
			FirstName: req.PaxContact.FirstName,
			LastName:  req.PaxContact.LastName,
			Phone:     req.PaxContact.Phone1,
		}
		if req.PaxContact.EmailAddress != nil {
						contact.Email = *req.PaxContact.EmailAddress
					}

					if errDB := db.Where("phone = ?", contact.Phone).FirstOrCreate(&contact).Error; errDB != nil {
						log.Printf("[BookFlight] fallback contact error: %v", errDB)
					}

					// SIMPAN TRANSACTION (FAILED)
					tripType := "oneway"
					if len(req.Trips) > 1 {
						tripType = "roundtrip"
					}

					flightTx := models.FlightTransaction{
						ContactID:      contact.ID,
						VMCode:         "CGK001",
						Amount:         0,
						TripType:       tripType,
						Status:         "failed", // ⬅️ ini pembeda utama
						SourcePlatform: "mobile",
					}

					if errDB := db.Create(&flightTx).Error; errDB != nil {
						log.Printf("[BookFlight] fallback transaction error: %v", errDB)
					} else {
						log.Printf("[BookFlight] fallback transaction saved ID: %d", flightTx.ID)
					}

					// OPTIONAL: simpan raw request untuk retry nanti
					// db.Create(&models.FailedBookingLog{Payload: req})

					c.JSON(http.StatusBadGateway, utils.ErrorResponse("Book flight failed", err.Error()))
					return
				}

				// ── SIMPAN CONTACT ──
				contact := models.Contact{
					FirstName: req.PaxContact.FirstName,
					LastName:  req.PaxContact.LastName,
					Phone:     req.PaxContact.Phone1,
				}
				if req.PaxContact.EmailAddress != nil {
					contact.Email = *req.PaxContact.EmailAddress
				}
				if err := db.Where("phone = ?", contact.Phone).FirstOrCreate(&contact).Error; err != nil {
					log.Printf("[BookFlight] Gagal simpan contact: %v", err)
					c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan contact", err.Error()))
					return
				}
				log.Printf("[BookFlight] Contact ID: %d", contact.ID)

				// ── SIMPAN TRANSACTION ──
				tripType := "oneway"
				if len(req.Trips) > 1 {
					tripType = "roundtrip"
				}

				totalAmount := result.Departure.TotalFare
				if result.Return.TotalFare > 0 {
					totalAmount += result.Return.TotalFare
				}

				flightTx := models.FlightTransaction{
					ContactID:      contact.ID,
					VMCode:         "CGK001",
					Amount:         totalAmount,
					TripType:       tripType,
				Status:         "pending",
			SourcePlatform: "mobile",
		}
		if err := db.Create(&flightTx).Error; err != nil {
			log.Printf("[BookFlight] Gagal simpan FlightTransaction: %v", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan transaksi", err.Error()))
			return
		}
		log.Printf("[BookFlight] FlightTransaction ID: %d", flightTx.ID)

		// ── SIMPAN BOOKING DEPARTURE ──
		depBooking := models.FlightBooking{
			TransactionID: flightTx.ID,
			PNRID:         result.Departure.PNRID,
			NTSA:          result.Departure.NTSA,
			Commission:    result.Departure.Commission,
			PublishFare:   result.Departure.TotalFare,
			Status:        "booked",
		}
		if err := db.Create(&depBooking).Error; err != nil {
			log.Printf("[BookFlight] Gagal simpan FlightBooking departure: %v", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan booking departure", err.Error()))
			return
		}
		log.Printf("[BookFlight] FlightBooking departure ID: %d, PNR: %s", depBooking.ID, depBooking.PNRID)

		// ── SIMPAN BOOKING RETURN (roundtrip beda maskapai = 2 PNR) ──
		var retBooking *models.FlightBooking
		if result.Return.PNRID != "" {
			rb := models.FlightBooking{
				TransactionID: flightTx.ID,
				PNRID:         result.Return.PNRID,
				NTSA:          result.Return.NTSA,
				Commission:    result.Return.Commission,
				PublishFare:   result.Return.TotalFare,
				Status:        "booked",
			}
			if err := db.Create(&rb).Error; err != nil {
				log.Printf("[BookFlight] Gagal simpan FlightBooking return: %v", err)
				c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan booking return", err.Error()))
				return
			}
			log.Printf("[BookFlight] FlightBooking return ID: %d, PNR: %s", rb.ID, rb.PNRID)
			retBooking = &rb
		}

		// ── SIMPAN FLIGHT SEGMENTS ke DB ──
		for _, t := range req.Trips {
			bookingID := depBooking.ID
			if t.TripType == "return" && retBooking != nil {
				bookingID = retBooking.ID
			}
			for _, seg := range t.Segments {
				departDateParsed, _ := time.Parse("2006-01-02", seg.DepartDate)
				db.Create(&models.FlightSegment{
					BookingID:   bookingID,
					AirlineCode: t.AirlineCode,
					FlightNo:    seg.FlightNo,
					FromCity:    seg.FromCity,
					ToCity:      seg.ToCity,
					ClassCode:   seg.ClassCode,
					DepartTime:  seg.DepartTime,
					ArriveTime:  seg.ArriveTime,
					DepartDate:  departDateParsed.Format("02-Jan-2006"),
				})
			}
		}

		// ── SIMPAN PASSENGERS ──
		for i, p := range req.ListPax {
			paxID := i + 1

			var dobPtr *time.Time
			if p.DOB != nil && *p.DOB != "" {
				t, err := time.Parse("2006-01-02", *p.DOB)
				if err == nil {
					dobPtr = &t
				} else {
					log.Printf("[BookFlight] listpax[%d]: gagal parse DOB '%s': %v", i, *p.DOB, err)
				}
			}

			var passportExpPtr *time.Time
			var passportNo, nationality, issuing string
			if p.Passport != nil {
				passportNo = p.Passport.No
				nationality = p.Passport.Nationality
				issuing = p.Passport.IssuingCountry
				if p.Passport.Expiry != "" {
					t, err := time.Parse("2006-01-02", p.Passport.Expiry)
					if err == nil {
						passportExpPtr = &t
					} else {
						log.Printf("[BookFlight] listpax[%d]: gagal parse passport expiry '%s': %v", i, p.Passport.Expiry, err)
					}
				}
			}

			lastName := p.LastName
			if lastName == "" {
				lastName = "/"
			}

			pax := models.FlightPassenger{
				BookingID:              depBooking.ID,
				Title:                  p.Title,
				FullName:               p.FirstName + " " + lastName,
				Type:                   p.Type,
				PaxID:                  paxID,
				DOB:                    dobPtr,
				AdultAssoc:             p.AdultAssoc,
				PassportNo:             passportNo,
				PassportNationality:    nationality,
				PassportIssuingCountry: issuing,
				PassportExpiry:         passportExpPtr,
			}
			if err := db.Create(&pax).Error; err != nil {
				log.Printf("[BookFlight] Gagal simpan FlightPassenger[%d] '%s': %v", i, pax.FullName, err)
				c.JSON(http.StatusInternalServerError, utils.ErrorResponse(
					fmt.Sprintf("Gagal simpan penumpang %s", pax.FullName), err.Error(),
				))
				return
			}
			log.Printf("[BookFlight] FlightPassenger saved ID: %d, Name: %s", pax.ID, pax.FullName)
		}

		// ── RESPONSE ──
		resp := gin.H{
			"transaction_id": flightTx.ID,
			"trip_type":      tripType,
			"total_amount":   totalAmount,
			"departure": gin.H{
				"booking_id": depBooking.ID,
				"pnr":        depBooking.PNRID,
				"ntsa":       depBooking.NTSA,
				"commission": depBooking.Commission,
				"total_fare": depBooking.PublishFare,
				"status":     depBooking.Status,
			},
		}
		if retBooking != nil {
			resp["return"] = gin.H{
				"booking_id": retBooking.ID,
				"pnr":        retBooking.PNRID,
				"ntsa":       retBooking.NTSA,
				"commission": retBooking.Commission,
				"total_fare": retBooking.PublishFare,
				"status":     retBooking.Status,
			}
		}

		c.JSON(http.StatusCreated, utils.SuccessResponse("Flight booked successfully", resp))
	}
}

// POST /api/v1/voltras/flight/pnr/retrieve
func RetrievePNRHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RetrievePNRRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		result, err := volClient.RetrievePNR(req.PNR)
		if err != nil {
			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Retrieve PNR failed", err.Error()))
			return
		}

		if err := db.Model(&models.FlightBooking{}).
			Where("pnrid = ?", req.PNR).
			Updates(map[string]interface{}{
				"status":       strings.ToLower(result.Status),
				"booking_code": result.BookingCode,
			}).Error; err != nil {
			log.Printf("[RetrievePNR] Gagal update FlightBooking PNR %s: %v", req.PNR, err)
		}

		if strings.EqualFold(result.Status, "TICKETED") {
			for _, pax := range result.ListPax {
				if pax.TicketNo != "" {
					if err := db.Model(&models.FlightPassenger{}).
						Joins("JOIN t_flight_bookings ON t_flight_bookings.id = t_flight_passengers.booking_id").
						Where("t_flight_bookings.pnrid = ? AND t_flight_passengers.pax_id = ?", req.PNR, pax.PaxNo).
						Update("ticket_number", pax.TicketNo).Error; err != nil {
						log.Printf("[RetrievePNR] Gagal update ticket_number pax %s PNR %s: %v", pax.PaxNo, req.PNR, err)
					}
				}
			}
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("PNR retrieved", result))
	}
}

// POST /api/v1/voltras/flight/ticket
func TicketingFlightHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TicketingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		var booking models.FlightBooking
		if err := db.First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}

		result, err := volClient.TicketingFlight(booking.PNRID)
		if err != nil {
			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Ticketing failed", err.Error()))
			return
		}

		if err := db.Model(&booking).Update("status", "ticketed").Error; err != nil {
			log.Printf("[TicketingFlight] Gagal update status booking ID %d: %v", booking.ID, err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal update status booking", err.Error()))
			return
		}
		log.Printf("[TicketingFlight] Booking ID %d status updated to ticketed", booking.ID)

		c.JSON(http.StatusOK, utils.SuccessResponse("Ticketing successful", gin.H{
			"booking_id": booking.ID,
			"pnr":        result.PNRID,
			"status":     "ticketed",
		}))
	}
}

// POST /api/v1/voltras/flight/autoticket
func AutoTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AutoTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		var booking models.FlightBooking
		if err := db.First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}

		if booking.Status != "booked" {
			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse(
				"Booking status tidak valid untuk auto ticket",
				fmt.Sprintf("status saat ini: %s", booking.Status),
			))
			return
		}

		result, err := volClient.AutoTicket(booking.PNRID, req.Email, req.SendTicket)
		if err != nil {
			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Auto ticket failed", err.Error()))
			return
		}

		remark := models.FlightRemark{
			BookingID: booking.ID,
			Remark:    fmt.Sprintf("AUTOTICKET VA: %s | TIMELIMIT: %s", result.AccountNo, result.TimeLimit),
		}
		if err := db.Create(&remark).Error; err != nil {
			log.Printf("[AutoTicket] Gagal simpan FlightRemark booking ID %d: %v", booking.ID, err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan remark", err.Error()))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Auto ticket created", gin.H{
			"booking_id": booking.ID,
			"pnr":        result.PNRID,
			"account_no": result.AccountNo,
			"time_limit": result.TimeLimit,
		}))
	}
}

// POST /api/v1/voltras/flight/ticket/cancel
func CancelTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CancelTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		var booking models.FlightBooking
		if err := db.First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}

		if booking.Status == "cancelled" {
			c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse("Booking sudah dibatalkan", ""))
			return
		}

		result, err := volClient.CancelTicket(booking.PNRID)
		if err != nil {
			c.JSON(http.StatusBadGateway, utils.ErrorResponse("Cancel failed", err.Error()))
			return
		}

		if err := db.Model(&booking).Update("status", "cancelled").Error; err != nil {
			log.Printf("[CancelTicket] Gagal update status booking ID %d: %v", booking.ID, err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal update status booking", err.Error()))
			return
		}

		remarkText := "CANCELLED"
		if req.Reason != "" {
			remarkText = "CANCELLED: " + req.Reason
		}
		remark := models.FlightRemark{BookingID: booking.ID, Remark: remarkText}
		if err := db.Create(&remark).Error; err != nil {
			log.Printf("[CancelTicket] Gagal simpan remark booking ID %d: %v", booking.ID, err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal simpan remark cancel", err.Error()))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Ticket cancelled successfully", gin.H{
			"booking_id": booking.ID,
			"pnr":        result.PNRID,
			"status":     "cancelled",
		}))
	}
}

// GET /api/v1/voltras/flight/print?pnr=XXX
func PrintTicketHandler(c *gin.Context) {
	pnr := c.Query("pnr")
	if pnr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("query param 'pnr' wajib diisi", ""))
		return
	}
	raw, contentType, err := volClient.PrintTicket(pnr)
	if err != nil {
		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Print ticket failed", err.Error()))
		return
	}
	filename := "ticket_" + pnr + ".html"
	if strings.Contains(contentType, "pdf") {
		filename = "ticket_" + pnr + ".pdf"
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, contentType, raw)
}

// GET /api/v1/voltras/flight/print-insurance?pnr=XXX
func PrintInsuranceHandler(c *gin.Context) {
	pnr := c.Query("pnr")
	if pnr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("query param 'pnr' wajib diisi", ""))
		return
	}
	raw, err := volClient.PrintInsurance(pnr)
	if err != nil {
		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Print insurance failed", err.Error()))
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="insurance_%s.pdf"`, pnr))
	c.Data(http.StatusOK, "application/pdf", raw)
}

// POST /api/v1/voltras/flight/advance-retrieve
func AdvanceRetrieveHandler(c *gin.Context) {
	var req models.AdvanceRetrieveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	if req.Status == "" {
		req.Status = "ALL"
	}
	validStatus := map[string]bool{"ALL": true, "BOOKED": true, "TICKETED": true, "CANCELED": true}
	if !validStatus[strings.ToUpper(req.Status)] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(
			"status tidak valid, gunakan: ALL, BOOKED, TICKETED, CANCELED", "",
		))
		return
	}

	fromDate, err := formatVolDate(req.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("from_date tidak valid (YYYY-MM-DD)", err.Error()))
		return
	}
	toDate, err := formatVolDate(req.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("to_date tidak valid (YYYY-MM-DD)", err.Error()))
		return
	}

	result, err := volClient.AdvanceRetrieve(req.Airline, fromDate, toDate, strings.ToUpper(req.Status))
	if err != nil {
		c.JSON(http.StatusBadGateway, utils.ErrorResponse("Advance retrieve failed", err.Error()))
		return
	}

	if result.ListPNR == nil {
		result.ListPNR = []AdvanceRetrievePNR{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Advance retrieve success", gin.H{
		"total": len(result.ListPNR),
		"list":  result.ListPNR,
	}))
}