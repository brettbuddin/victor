package victor

import (
    "log"
    "github.com/brettbuddin/victor/utils/google"
)

func (self *Robot) registerDefaultAbilities() {
    self.Respond("ping", func(msg *TextMessage) {
        msg.Send("Pong!")
    })

    self.Respond("(image|img) (.*)", func(msg *TextMessage) {
        result, err := google.ImageSearch(msg.Matches()[3])

        if err != nil {
            log.Print(err)
            return
        }

        if result == "" {
            msg.Send("I didn't find anything.")
            return
        }
        
        msg.Send(result)    
    })

    self.Respond("show users", func(msg *TextMessage) {
        list := ""
        for _, user := range self.users {
            list += user.Name + "\n" 
        }

        msg.Send(list)
    })
}
