/* Simple analysis on a monophonic single frequency WAV file
 * Returns a channel of type Pitch.
 * The type Pitch carries the Hop number at which the analysis is made, the pitch detected at that Hop and the confidence of the detection.
 * pitch of -1 is returned in case of silence or failure to detect the time period by the algorithm.
 * Pitch is analysed for every chunk of size defined by the third parameter, the hopSize
 * The hopSize is time-frequency trade-off. For better time resolution, use a smaller chunk size. 
 * However, this would limit the detection of lower frequencies.
 * 
 * Theoretically, you would want 2*SamplingRate/FrequencyToBeDetected number of samples for analysis.
 * So, with 2048 chunk size, you are looking at about 43 Hz lowest possible frequency.
 */


package main 
import "github.com/mrnikho/yingo"
import "fmt"

func main () {
	
	pitchChannel := yingo.MonoAnalyser("piano.wav", true, 2048)
	
	for pitch := range pitchChannel {
		fmt.Println(pitch)
	}
}
