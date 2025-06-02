package server

// type IEMConfig struct {
// 	StartDate      time.Time
// 	EndDate        *time.Time
// 	Product        string
// 	Office         string
// 	MaxConcurrency int
// }

// type ProgressReader struct {
// 	Reader       io.Reader
// 	Total        int64
// 	Downloaded   int64
// 	LastReported time.Time
// }

// // Read tracks the progress as the file downloads
// func (pr *ProgressReader) Read(p []byte) (int, error) {
// 	n, err := pr.Reader.Read(p)
// 	pr.Downloaded += int64(n)

// 	// Report progress every 500ms to avoid spamming the console
// 	if time.Since(pr.LastReported) > 500*time.Millisecond {
// 		fmt.Printf("Downloaded: %d/%d bytes (%.2f%%)\n",
// 			pr.Downloaded, pr.Total, float64(pr.Downloaded)/float64(pr.Total)*100)
// 		pr.LastReported = time.Now()
// 	}

// 	return n, err
// }

// func IEM(config IEMConfig, minLog int) {
// 	offices := []string{}

// 	sconfig := ServerConfig{
// 		MinLog: minLog,
// 	}

// 	server, err := New(sconfig)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// 	if config.Office == "ALL" {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		rows, err := server.DB.Query(ctx, `
// 		SELECT id FROM postgis.offices ORDER BY id`)
// 		if err != nil {
// 			slog.Error("failed to get offices: " + err.Error())
// 			return
// 		}
// 		defer rows.Close()

// 		for rows.Next() {
// 			var office string
// 			err := rows.Scan(&office)
// 			if err != nil {
// 				slog.Error("failed to scan office: " + err.Error())
// 				return
// 			}
// 			offices = append(offices, office)
// 		}
// 		if rows.Err() != nil {
// 			slog.Error("failed to get offices: " + rows.Err().Error())
// 			return
// 		}
// 		if len(offices) == 0 {
// 			slog.Error("no rows returned for ALL")
// 			return
// 		}
// 	} else {
// 		offices = append(offices, config.Office)
// 	}

// 	startDate := config.StartDate
// 	endDate := config.StartDate.Add(24 * time.Hour)
// 	if config.EndDate != nil {
// 		endDate = *config.EndDate
// 	}
// 	if startDate.After(endDate) {
// 		s := startDate
// 		startDate = endDate
// 		endDate = s
// 	}

// 	for d := startDate; !d.Equal(endDate); d = d.Add(24 * time.Hour) {

// 		start := d.Format("2006-01-02")
// 		end := d.Add(24 * time.Hour).Format("2006-01-02")
// 		data, err := retrieveIEMData(start, end, offices, config.Product)
// 		if err != nil {
// 			slog.Error("failed retrieving data from IEM: "+err.Error(), "startDate", start, "endDate", end)
// 			return
// 		}

// 		for _, text := range data {

// 			h, err := handler.New(server.DB, server.MinLog)
// 			if err != nil {
// 				slog.Error("failed to create handler: " + err.Error())
// 				return
// 			}

// 			err = h.Handle(text, time.Now())
// 			if err != nil {
// 				slog.Error("failed to handle item: " + err.Error())
// 				continue
// 			}
// 		}
// 	}
// }

// func retrieveIEMData(start string, end string, offices []string, product string) ([]string, error) {

// 	// Forse zip format
// 	url := "https://mesonet.agron.iastate.edu/cgi-bin/afos/retrieve.py?fmt=zip&limit=9999&"
// 	for _, o := range offices {
// 		url += fmt.Sprintf("pil=%s%s&", product, o)
// 	}

// 	url += fmt.Sprintf("sdate=%s&edate=%s", start, end)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("received status %d while getting from IEM", resp.StatusCode)
// 	}

// 	totalSize := resp.ContentLength
// 	if totalSize <= 0 {
// 		slog.Warn("Content length unknown, progress will be inaccurate.")
// 	}

// 	progressReader := &ProgressReader{
// 		Reader:       resp.Body,
// 		Total:        totalSize,
// 		LastReported: time.Now(),
// 	}

// 	body, err := io.ReadAll(progressReader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read ZIP: %s", err.Error())
// 	}
// 	slog.Info("Download complete.")

// 	reader := bytes.NewReader(body)
// 	zipReader, err := zip.NewReader(reader, int64(len(body)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := []string{}

// 	for _, file := range zipReader.File {
// 		zippedFile, err := file.Open()
// 		if err != nil {
// 			return nil, err
// 		}

// 		content, err := io.ReadAll(zippedFile)
// 		if err != nil {
// 			zippedFile.Close()
// 			return nil, err
// 		}
// 		zippedFile.Close()

// 		data = append(data, string(content))
// 	}

// 	return data, nil
// }
