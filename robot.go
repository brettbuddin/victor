package victor

import (
    "regexp"
    "log"
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

func NewRobot(adapter string, options map[string]string) RobotAdapter {
    name := "robot"

    if _, exists := options["name"]; exists {
        name = options["name"]
    }

    robot := &Robot{
       name: name,
       options: options,
       listeners: make([]*Listener, 0, 1),
    }

    var aRobot RobotAdapter

    if adapter == "shell" {
        aRobot = &Shell{Robot: robot}
    } else if adapter == "campfire" {
        aRobot = &Campfire{Robot: robot}
    } else {
        panic("Unknown adapter.")
    }

    return aRobot
}

func (self *Robot) Hear(expStr string, callback func(*TextMessage)) {
    exp, _ := regexp.Compile(expStr)
    self.listeners = append(self.listeners, NewListener(exp, callback))
}

func (self *Robot) Respond(expStr string, callback func(*TextMessage)) {
    expWithNameStr := "(" + self.name + "[:,]?)\\s*(?:" + expStr + ")" 
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
