/* Credit goes to Henri B for the mentorship and Coding */

package main

import "fmt"
import "net/http"
import "io/ioutil"
import "strings"
import "math/rand"
import "time"

//num of threads for reading words
//more threads = faster output, more varied output
//keep alive < max
var maxraw int = 10
var aliveraw int = 0

func main() {
   //setup the output, and initiate the printer
   chResult := make(chan string)
   go printResults(chResult)
   
   //get initial set
   go getPastebinIDs(chResult)
   
   //keep the main thread for dieing, 300 is some large number
   for {
      time.Sleep(time.Duration(300) * time.Second)
   }
}

//output the results
func printResults(chResult chan string) {
   for {
      fmt.Print(<-chResult+" ")
   }
}

//Open pastebin.com to get a new set of raws to parse
func getPastebinIDs(chResult chan string) {
   //get the body of the site (comes as binary)
   response, _ := http.Get("https://pastebin.com")
   body, _ := ioutil.ReadAll(response.Body)

   //get the proper subsection (a listing of (8?) links to raws)
   ids := strings.Split(string(body), "Public Pastes")
   ids = strings.Split(ids[1], "iframe")
   
   //split so each raw id is at the begining of a new line, 0th is junk
   ids = strings.Split(ids[0], "href=\"/")
   
   //wait for a spot to open up, then begin parsing the raw
   //30 seconds is arbitrary
   for i := 1; i<len(ids); i++ {
      for ; aliveraw > maxraw; {
         time.Sleep(time.Duration(30) * time.Second)
      }
      go getPastebinRaw(ids[i][:8], chResult)
   }
   
   //close the connection (prob not important), wait a bit, repeat
   response.Body.Close()
   time.Sleep(time.Duration(30) * time.Second)
   go getPastebinIDs(chResult)
}

//append each word of the raw to the output, waiting a random time between each 
func getPastebinRaw(id string, chResult chan string) {
   //take a spot in the raws, then get the page
   aliveraw++
   response, _ := http.Get("https://pastebin.com/raw/"+id)
   body, _ := ioutil.ReadAll(response.Body)
   
   //split on space (could try "\r\n" (new lines) but space seems better
   split := strings.Split(string(body), " ")
   
   //if word is not too long, append it
   //then wait some random amount
   for i:=0; i<len(split); i++ {
      if len(string(split[i])) < 10 {
         chResult<-string(split[i])
      }
      for j := 0; j<20; j++ {
         time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
      }
   }
   
   //finished, return raw slot and close
   aliveraw--
   response.Body.Close()
}
