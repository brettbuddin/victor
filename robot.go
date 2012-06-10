package victor

import (
    "regexp"
    "log"

    "github.com/brettbuddin/victor/utils/google"
)

type RobotAdapter interface {
    Hear(string, func(*TextMessage))
    Respond(string, func(*TextMessage))
    Run()
}

type Robot struct {
    name      string
    options   map[string]string
    adapter   RobotAdapter
    listeners []*Listener
    users     []*User
}

func NewRobot(adapter string, name string) RobotAdapter {
    if name == "" {
        name = "victor"
    }

    robot := &Robot{
        name: name,
        listeners: make([]*Listener, 0, 1),
    }

    robot.registerDefaultAbilities()

    return robot.LoadAdapter(adapter)
}

func (self *Robot) LoadAdapter(name string) RobotAdapter {
    var adapter RobotAdapter

    if name == "shell" || name == "" {
        adapter = NewShell(self)
    } else if name == "campfire" {
        adapter = NewCampfire(self)
    } else {
        log.Panic("Unkown adapter.")
    }

    return adapter
}

func (self *Robot) Hear(expStr string, callback func(*TextMessage)) {
    exp, _ := regexp.Compile(expStr)
    self.listeners = append(self.listeners, NewListener(exp, callback))
}

func (self *Robot) Respond(expStr string, callback func(*TextMessage)) {
    expWithNameStr := "^(" + self.name + "[:,]?)\\s*(?:" + expStr + ")" 
    exp, _         := regexp.Compile(expWithNameStr)

    self.listeners = append(self.listeners, NewListener(exp, callback))
}

func (self *Robot) Receive(msg *TextMessage) {
    for _, listener := range self.listeners {
        if listener.Test(msg) {
            log.Printf("Listener /%s/ triggered by '%s'", listener.Exp.String(), msg.Body)
            listener.Callback(msg)
        }
    }
}

func (self *Robot) RememberUser(user *User) {
    self.users = append(self.users, user)
}

func (self *Robot) UserForId(id int) *User {
    for _, user := range self.users {
        if user.Id == id {
            return user
        }
    }

    return nil
}

func (self *Robot) Shutdown() {
    log.Print("See ya!")
}

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

    self.Respond("(list|show) users", func(msg *TextMessage) {
        list := ""
        for _, user := range self.users {
            list += user.Name + "\n" 
        }

        msg.Paste(list)
    })
}
