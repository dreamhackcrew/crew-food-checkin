package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"os"
	"os/exec"
	"time"
	"io"
	//"strings"
)


// Globals
var led [2]chan byte
var port io.ReadWriteCloser
var port_ctrl chan string

func main() {

    // Startup - IO
    led[0] = make(chan byte,10);
    led[1] = make(chan byte,10);
    go setupGPIO();

    // Startup - serial port
    go startSerialPort();

    // Start main loop
    running := true
    ticker := time.NewTicker(time.Second / 10) // 5 Hz
    for running {
        select {
        case <-ticker.C:
	}
    }
}


func check( barcode string ) {
    cmd := exec.Command("/home/pi/crew-food-checkin/scripts/check.php",barcode)
    stdout, err := cmd.Output(); 
    
    if err != nil { log.Println(err) }

    packets := string(stdout[:]);
    log.Print("-----------------");
    log.Print(packets)
    log.Print("-----------------");
    log.Print(packets)
    switch packets {
	case "success":
		led[0] <- 50
	case "fail":
		led[1] <- 50
	default:
		led[0] <- 10
		led[1] <- 10
    }
}


// GPIO
func setupGPIO() {/*{{{*/
    cmd := exec.Command("gpio","export","17","out")
    if err := cmd.Run(); err != nil { log.Println("1");log.Fatalln(err) }
    cmd = exec.Command("gpio","export","18","out")
    if err := cmd.Run(); err != nil { log.Println("2"); log.Fatalln(err) }

    go led_runner(0);
    go led_runner(1);
}/*}}}*/
func led_runner( index byte ) {/*{{{*/
    var f *os.File
    var err error
    switch index {
        case 0:
            f, err = os.OpenFile("/sys/class/gpio/gpio17/value", os.O_WRONLY, 0222)
            if err != nil {
                log.Print(err)
                return
            }
            defer f.Close();
            defer f.Write([]byte("0"))
            f.Write([]byte("1"))
        default:
            f, err = os.OpenFile("/sys/class/gpio/gpio18/value", os.O_WRONLY, 0222)
            if err != nil {
                log.Print(err)
                return
            }
            defer f.Close();
            defer f.Write([]byte("0"))
            f.Write([]byte("1"))
    }

    ticker := time.NewTicker(time.Second / 50) // 5 Hz
    running := true;

    status := "0";
    var value byte;
    var cnt byte;
    value = 255;
    for running {
        select {
            case <-ticker.C:
                if value > 0 && value < 255 {
                    if cnt > value {
                        value = 0;
                        status = "1";
                        f.Write([]byte(status))
                    }

                    cnt++;
                } 
            case data := <-led[index]:
		cnt = 0;
                value = data
                if value == 0 {
                    status = "1";
                    f.Write([]byte(status))
                } else {
                    status = "0";
                    f.Write([]byte(status))
                }
        }
    }
}/*}}}*/

// Serial port
func startSerialPort() {/*{{{*/
    // Connect to the serial port
    //connectPort("/dev/ttyUSB0");
    connectPort("/dev/ttyACM0");

    for {
        time.Sleep(time.Second/80) // 50 Hz
        readPort();
    }
}/*}}}*/
func connectPort( s string ) {/*{{{*/
    //c := &serial.Config{Name: s, Baud: 115200}
    c := &serial.Config{Name: s, Baud: 19200}
    var err error
    port, err = serial.OpenPort(c)
    if err != nil {
            log.Fatal(err)
    }
}/*}}}*/
func readPort() {/*{{{*/
    buf := make([]byte, 4096)
    n, err := port.Read(buf)
    if err != nil {
        return;
        log.Print(err)
        log.Fatal("Failed to connect to PORT")
    }
    packets := string(buf[:n]);
    log.Print(packets)

    go check(packets);

}/*}}}*/
func writePort(s string) {/*{{{*/
    fmt.Print("Set addr:");
    fmt.Print( []byte(s) );
    fmt.Print("\n");
    _,err := port.Write( []byte(s) )
    if err != nil {
            log.Fatal(err)
    }

    /*packets := strings.Split(s,string(0x03))

    for _,packet := range packets {
        start := strings.Index(packet,string(0x02));

        if start == -1 {
            continue;
        }

        packet = packet[start+1:]; // Remove the start char

        if len(packet) != 3 {
            continue;
        }

        select {
            case dmx <- []byte(packet):
            default:
        }
    }*/

}/*}}}*/
