package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "sync"
    "time"
)

var WorkGroup sync.WaitGroup;

func ScanPort(port int, Target string, Timeout int) string {
    
    Time, _ := time.ParseDuration(os.Args[5] + "s")
    
    conn, err := net.DialTimeout("tcp", Target+":"+strconv.Itoa(port), Time )

    if err != nil {
        fmt.Println("ERR::" + strconv.Itoa(port) + ">" + err.Error())
        return "Port " + strconv.Itoa(port) + " failed."
    } else {
        recvBuf := make([]byte, 1024)
        _, err = conn.Read(recvBuf[:])
        
        if err != nil {
    
            return err.Error()
            
        } else {
            
            if string(recvBuf[:]) == "" {
                fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
                _, err = conn.Read(recvBuf[:])
            }
            if string(recvBuf[:]) == "" {
                fmt.Fprintf(conn, "GET / HTTP/1.1\r\n\r\n")
                _, err = conn.Read(recvBuf[:])
            }
            if string(recvBuf[:]) == "" {
                fmt.Fprintf(conn, "GET / HTTP/2.0\r\n\r\n")
                _, err = conn.Read(recvBuf[:])
            }
            
            conn.Close()
                
            return "##################### Port "+strconv.Itoa(port)+" ###################\n\n" + string(recvBuf[:]) + "\n\n###############################################################"
		}

	}
}

func EnumeratePorts(Start int, End int, Target string, Timeout int) {
    var Banners = make([]string, End-Start)
    
    for i := Start; i < End; i++ {
        Banners[i-Start] = ScanPort(i, Target, Timeout)
    } 
    
    for i := 0; i < End-Start; i++ {
        fmt.Println(Banners[i])
    }
    
    WorkGroup.Done()
}

func main() {
    
    if os.Args[1] == "--help" {
        fmt.Println("Useage: gopscan <Target IP> <# Threads> <Start Port> <End Port> <Timeout (seconds)>")
        return
    }
    
    Target := os.Args[1]
    Threads, _ := strconv.Atoi(os.Args[2])
    StartPort, _ := strconv.Atoi(os.Args[3])
    EndPort, _ := strconv.Atoi(os.Args[4])
    Timeout, _ := strconv.Atoi(os.Args[5])

    var TotalPorts = EndPort - StartPort
    var PortsPerThread = TotalPorts / Threads
    var RemainingPorts = TotalPorts % Threads
    
    fmt.Println("################ CONFIG ##############\n")
    fmt.Println(" >Target = " + Target + "\n")
    fmt.Println(" >Start Port = " + strconv.Itoa(StartPort) + "\n")
    fmt.Println(" >End Port = " + strconv.Itoa(EndPort) + "\n")
    fmt.Println(" >Number of Threads = " + strconv.Itoa(Threads) + "\n")
    fmt.Println(" >Connection Timeout = " + strconv.Itoa(Timeout) + "\n")
    
    for i := 1; i <= Threads; i++ {
        fmt.Println(" >Thread"+ strconv.Itoa(i) + " : {" + strconv.Itoa( ((i-1)*PortsPerThread) + StartPort ) +":"+ strconv.Itoa( (i)*PortsPerThread + StartPort ) + "}\n")
    }
    
    fmt.Println("######################################\n")

    WorkGroup.Add(Threads)
    
    for i := 0; i < Threads; i++ {
        if i == Threads-1 {
            go EnumeratePorts( (i)*PortsPerThread+StartPort , ((i+1)*PortsPerThread)+StartPort + RemainingPorts , Target, Timeout )
        } else {
            go EnumeratePorts( (i)*PortsPerThread+StartPort , (i+1)*PortsPerThread+StartPort , Target, Timeout )
        }
    }
    
    WorkGroup.Wait()
}
