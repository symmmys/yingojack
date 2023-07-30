package yingo

import "fmt"
import "github.com/xthexder/go-jack"
//type Mic struct {
	//chunkSize int
	//bufferSize int
	//bufferoptimize bool
//}var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port


func MicInput(chunkSize int, bfropt bool, pch *chan float32) {
	//yin variables
	//m = Mic{chunkSize: chnkSze, bufferSize: 100, bufferoptimize: bfropt }
	threshold := float32(0.05)
	bufferSize := 100
	//var pitch float32
	var bufferincrement int
	//pch := make(chan float32)
	
	if bfropt {
		bufferincrement = 1
	} else {
		bufferincrement = 100
	}
	
	fmt.Println("Enter loop")
	go func () {
	
	//portaudio.Initialize()
	//defer portaudio.Terminate()
	
	input := make([]int16, chunkSize)
	//stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(input), input)
	//chkErr(err)
	//defer stream.Close()
	
	//chkErr(stream.Start()
        client, status := jack.ClientOpen("go-Monitor", jack.NoStartServer)
        if status != 0 {
        fmt.Println("Status:", jack.StrError(status))
        return
        }
        defer client.Close()

        for i := 0; i < channels; i++ {
        portIn := client.PortRegister(fmt.Sprintf("in_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
        PortsIn = append(PortsIn, portIn)
        }
        for i := 0; i < channels; i++ {
        portOut := client.PortRegister(fmt.Sprintf("out_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
        PortsOut = append(PortsOut, portOut)
        }

        MonBuffer := []int16
        signalChan := make(chan []int16)
        process := func(nframes uint32) int {
            for i, in := range PortsIn {

                items := make([]int16, 0)
                samplesIn := in.GetBuffer(nframes)
                samplesOut := PortsOut[i].GetBuffer(nframes)
                for i2, sample := range samplesIn {
                    samplesOut[i2] = sample
                    items = append(items, int16(sample))
                }
                signalChan <- items
            }
            return 0
        }
        if code := client.SetProcessCallback(process); code != 0 {
        fmt.Println("Failed to set process callback:", jack.StrError(code))
        return
        }
        shutdown := make(chan struct{})
        client.OnShutdown(func() {
        fmt.Println("Shutting down")
        close(shutdown)
        })

        if code := client.Activate(); code != 0 {
        fmt.Println("Failed to activate client:", jack.StrError(code))
        return
        }

        fmt.Println(client.GetName())
		for {
			yin:= Yin{}
            MonBuffer = <- signalChan
            fmt.Println("MonBuffer: " , MonBuffer)
			fmt.Println("New Chunk")
			fmt.Println(MonBuffer)
			figurePitch(&yin, bufferincrement, input, pch, bufferSize, threshold)

		}
	}()
    go func() {
        for {
            pitchFound := <-*pch
            fmt.Println("pitchFound: " , pitchFound)
        }

    }()
        <-shutdown
}

func figurePitch(yin *Yin, bufferincrement int, input []int16, pch *chan float32, bufferSize int, threshold float32){
	var pitch float32
	fmt.Println("Processing")
	for pitch < 10 {
		//fmt.Println(bufferSize)
		if bufferSize >= len(input) {
			fmt.Println("Break")
			pitch = -1
			break
		}
		yin.YinInit(bufferSize, threshold)
		pitch = yin.GetPitch(&input)
		bufferSize += bufferincrement
	}
	//fmt.Println(bufferSize)
	fmt.Println(pitch)
	*pch<- pitch
}

func chkErr(err error){
	if err !=nil{
		panic(err)
	}
}
