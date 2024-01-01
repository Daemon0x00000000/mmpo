package main

/*
func GetPrinter(ctx *gousb.Context) (*gousb.Device, error) {

	// OpenDevices is used to find the devices to open.
	dev, err := ctx.OpenDeviceWithVIDPID(0x0483, 0x5740)

	// OpenDevices can occasionally fail, so be sure to check its return value.
	if err != nil || dev == nil {
		return nil, err
	}

	return dev, nil
}

func logError(err error) {
	log.Fatalf("Error : %v", err)
}

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	printer, err := GetPrinter(ctx)

	if err != nil {
		logError(err)
	}

	if printer != nil {
		defer printer.Close()
		err = printer.SetAutoDetach(true)
		if err != nil {
			log.Fatalf("%s.SetAutoDetach(true): %v", printer, err)
		}
		fmt.Printf("Infos : %s \n", usbid.Describe(printer.Desc))

		cfg, err := printer.Config(1)
		if err != nil {
			logError(err)
		}
		defer cfg.Close()

		intf, err := cfg.Interface(1, 0)
		if err != nil {
			logError(err)
		}
		defer intf.Close()

		fmt.Println(intf.Setting)

		out, err := intf.OutEndpoint(0x01)
		if err != nil {
			logError(err)
		}

		buf := []byte("M105\n\r")

		write, err := out.Write(buf)
		if write != 4 {
			log.Fatalf("%s.Write([5]): only %d bytes written, returned error is %v", out, write, err)
		}

	}

	bytes := []byte("M-179")

	fmt.Println(bytes)

}

*/

var printer Printer

func main() {

	var err error
	printer, err = NewPrinter()
	if err != nil {
		return
	}
	defer printer.Close()

	/*
		commands := [...]string{"M155 S0"}
		// "G0 Y-40", "G0 Y40", "G0 X-30", "G0 X30"

		for _, command := range commands {
			answer, ok := printer.Send(command)
			fmt.Printf("%s : %t\n", answer, ok)
		}

	*/

	printer.AutoTemperaturesJob()
	printer.Print("test")

	/*
		commands := [...]string{"M105", "M105"}
		// "G0 Y-40", "G0 Y40", "G0 X-30", "G0 X30"

		for i := 0; i < 50; i++ {

			for _, command := range commands {
				answer, ok := printer.Send(command)
				fmt.Printf("%s : %t\n", answer, ok)
			}
			time.Sleep(200 * time.Millisecond)
		}

		for {

		}

	*/
}
