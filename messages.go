// here be dragons
package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	CLIENT_MESSAGE = iota
	CLIENT_ERROR
	SERVER_MESSAGE
	SERVER_ERROR
	JOINPART_MESSAGE
	UPDATE_MESSAGE
	NOTICE_MESSAGE
	ACTION_MESSAGE
	PRIVATE_MESSAGE
	CUSTOM_MESSAGE
)

func msgTypeString(t int) string {
	switch t {
	case CLIENT_MESSAGE:
		return "CLIENT_MESSAGE"
	case CLIENT_ERROR:
		return "CLIENT_ERROR"
	case SERVER_MESSAGE:
		return "SERVER_MESSAGE"
	case SERVER_ERROR:
		return "SERVER_ERROR"
	case JOINPART_MESSAGE:
		return "JOINPART_MESSAGE"
	case UPDATE_MESSAGE:
		return "UPDATE_MESSAGE"
	case NOTICE_MESSAGE:
		return "NOTICE_MESSAGE"
	case ACTION_MESSAGE:
		return "ACTION_MESSAGE"
	case PRIVATE_MESSAGE:
		return "PRIVATE_MESSAGE"
	case CUSTOM_MESSAGE:
		return "CUSTOM_MESSAGE"
	}
	return "(unknown)"
}

func clientError(tab tabWithTextBuffer, msg ...string) {
	Println(CLIENT_ERROR, T(tab), msg...)
}
func clientMessage(tab tabWithTextBuffer, msg ...string) {
	Println(CLIENT_MESSAGE, T(tab), msg...)
}
func serverMessage(tab tabWithTextBuffer, msg ...string) {
	Println(SERVER_MESSAGE, T(tab), msg...)
}
func serverError(tab tabWithTextBuffer, msg ...string) {
	Println(SERVER_ERROR, T(tab), msg...)
}
func joinpartMessage(tab tabWithTextBuffer, msg ...string) {
	Println(JOINPART_MESSAGE, T(tab), msg...)
}
func updateMessage(tab tabWithTextBuffer, msg ...string) {
	Println(UPDATE_MESSAGE, T(tab), msg...)
}
func noticeMessage(tab tabWithTextBuffer, msg ...string) {
	Println(NOTICE_MESSAGE, T(tab), msg...)
}
func actionMessage(tab tabWithTextBuffer, msg ...string) {
	Println(ACTION_MESSAGE, T(tab), msg...)
}
func privateMessage(tab tabWithTextBuffer, msg ...string) {
	Println(PRIVATE_MESSAGE, T(tab), msg...)
}

type highlighterFn func(nick, msg string) bool

func noticeMessageWithHighlight(tab tabWithTextBuffer, hl highlighterFn, msg ...string) {
	PrintlnWithHighlight(NOTICE_MESSAGE, hl, T(tab), msg...)
}
func actionMessageWithHighlight(tab tabWithTextBuffer, hl highlighterFn, msg ...string) {
	PrintlnWithHighlight(ACTION_MESSAGE, hl, T(tab), msg...)
}
func privateMessageWithHighlight(tab tabWithTextBuffer, hl highlighterFn, msg ...string) {
	PrintlnWithHighlight(PRIVATE_MESSAGE, hl, T(tab), msg...)
}

func T(tabs ...tabWithTextBuffer) []tabWithTextBuffer { return tabs } // expected type, found ILLEGAL

func PrintlnWithHighlight(msgType int, hl highlighterFn, tabs []tabWithTextBuffer, msg ...string) {
	switch msgType {
	case NOTICE_MESSAGE:
		for _, tab := range tabs {
			logmsg := now() + " *** NOTICE: " + strings.Join(msg, " ")
			tab.Logln(logmsg)

			tab.Notify(true) // always put a * for NOTICE

			h := false
			if len(msg) >= 3 {
				h = hl(msg[1], strings.Join(msg[2:], " "))
			}
			if h && !mainWindowFocused {
				systray.ShowMessage("", logmsg)
			}
			tab.Println(parseString(noticeMsg(h, msg...)))
		}

	case PRIVATE_MESSAGE:
		nick, msg := msg[0], strings.Join(msg[1:], " ")
		logmsg := now() + " <" + nick + "> " + msg
		h := hl(nick, msg)
		if h && !mainWindowFocused {
			systray.ShowMessage("", logmsg)
		}
		for _, tab := range tabs {
			nick := colorNick(tab, h, nick)
			nick = leftpadNick(tab, nick)
			tab.Notify(h)
			tab.Logln(logmsg)
			tab.Println(parseString(privateMsg(h, nick, msg)))
		}

	case ACTION_MESSAGE:
		logmsg := now() + " *" + strings.Join(msg, " ") + "*"
		nick, msg := msg[0], strings.Join(msg[1:], " ")
		h := hl(nick, msg)
		if h && !mainWindowFocused {
			systray.ShowMessage("", logmsg)
		}
		for _, tab := range tabs {
			nick := colorNick(tab, h, nick)
			tab.Notify(h)
			tab.Logln(logmsg)
			tab.Println(parseString(actionMsg(h, nick, msg)))
		}
	default:
		log.Printf("highlighting unsupported for msgType %v", msgTypeString(msgType))
	}
}

func Println(msgType int, tabs []tabWithTextBuffer, msg ...string) {
	if len(msg) == 0 {
		log.Printf("tried to print an empty line of type %v", msgTypeString(msgType))
		return
	}

	switch msgType {
	case CLIENT_MESSAGE:
		for _, tab := range tabs {
			tab.Println(parseString(clientMsg(msg...)))
		}

	case CLIENT_ERROR:
		for _, tab := range tabs {
			tab.Errorln(parseString(clientErrorMsg(msg...)))
		}

	case SERVER_MESSAGE:
		text, styles := parseString(serverMsg(msg...))
		for _, tab := range tabs {
			tab.Logln(text)
			tab.Println(text, styles)
		}

	case SERVER_ERROR:
		text, styles := parseString(serverErrorMsg(msg...))
		for _, tab := range tabs {
			tab.Logln(text)
			tab.Errorln(text, styles)
		}

	case JOINPART_MESSAGE:
		if !clientCfg.HideJoinParts {
			text, styles := parseString(joinpartMsg(msg...))
			for _, tab := range tabs {
				tab.Logln(text)
				tab.Println(text, styles)
			}
		}

	case UPDATE_MESSAGE:
		// TODO(tso): option to hide?
		text, styles := parseString(updateMsg(msg...))
		for _, tab := range tabs {
			tab.Logln(text)
			tab.Println(text, styles)
		}

	case NOTICE_MESSAGE:
		for _, tab := range tabs {
			tab.Notify(true)
			tab.Logln(now() + " *** NOTICE: " + strings.Join(msg, " "))
			tab.Println(parseString(noticeMsg(false, msg...)))
		}

	case PRIVATE_MESSAGE:
		nick, msg := msg[0], strings.Join(msg[1:], " ")
		logmsg := now() + " <" + nick + "> " + msg
		for _, tab := range tabs {
			nick := colorNick(tab, false, nick)
			nick = leftpadNick(tab, nick)
			tab.Notify(false)
			tab.Logln(logmsg)
			tab.Println(parseString(privateMsg(false, nick, msg)))
		}

	case ACTION_MESSAGE:
		logmsg := now() + " *" + strings.Join(msg, " ") + "*"
		nick, msg := msg[0], strings.Join(msg[1:], " ")
		for _, tab := range tabs {
			nick := colorNick(tab, false, nick)
			tab.Notify(false)
			tab.Logln(logmsg)
			tab.Println(parseString(actionMsg(false, nick, msg)))
		}

	case CUSTOM_MESSAGE:
		text, styles := parseString(strings.Join(msg, " "))
		for _, tab := range tabs {
			tab.Println(text, styles)
		}

	default:
		log.Printf("Println: unhandled msgType: %v", msgType)
	}
}

func clientMsg(text ...string) string {
	return color(strings.Join(text, " "), DarkGrey)
}

func clientErrorMsg(text ...string) string {
	return color(strings.Join(text, " "), Red)
}

func serverErrorMsg(text ...string) string {
	if len(text) < 2 {
		return fmt.Sprintf("wrong argument count for server error: want 2 got %d:\n%#v",
			len(text), text)
	}
	return color(now(), Red) +
		color("E", White, Red) +
		color("("+text[0]+"): "+strings.Join(text[1:], " "), Red)
}

func serverMsg(text ...string) string {
	if len(text) < 2 {
		return fmt.Sprintf("wrong argument count for server message: want 2 got %d:\n%#v",
			len(text), text)
	}
	return color(now()+"S("+text[0]+"): "+strings.Join(text[1:], " "), DarkGray)
}

func joinpartMsg(text ...string) string {
	return color(now(), LightGray) + " " + italic(color(strings.Join(text, " "), Orange))
}

func updateMsg(text ...string) string {
	return color(now(), DarkGray) + " " + color(strings.Join(text, " "), DarkGrey)
}

func noticeMsg(hl bool, text ...string) string {
	if len(text) < 3 {
		return fmt.Sprintf("wrong argument count for notice: want 3, got %d:\n%v", len(text), text)
	}
	line := color(now(), LightGray) +
		color("N", White, Orange) +
		color("("+text[0]+"->"+text[1]+"):", Orange)
	if hl {
		line += "»"
		line += color(strings.Join(text[2:], " "), White, Orange)
	} else {
		line += " "
		line += strings.Join(text[2:], " ")
	}
	return line
}

func actionMsg(hl bool, text ...string) string {
	line := color(now(), LightGray)
	if hl {
		line += "»"
	} else {
		line += " "
	}
	return line + "*" + strings.TrimSpace(strings.Join(text, " ")) + "*"
}

func privateMsg(hl bool, text ...string) string {
	if len(text) < 2 {
		return fmt.Sprintf("wrong argument count for notice: want 2, got %d:\n%v", len(text), text)
	}
	nick := text[0]
	line := color(now(), LightGray)
	if hl {
		line += "»"
	} else {
		line += " "
	}
	return line + nick + " " + strings.Join(text[1:], " ")

}

func colorNick(tab tabWithTextBuffer, hl bool, nick string) string {
	if hl {
		bg := tab.NickColor(nick)
		fg := White
		if !colorVisible(0xffffff, colorPalette[bg]) {
			fg = Black
		}
		return color(nick, fg, bg)
	}
	return color(nick, tab.NickColor(nick))
}

func leftpadNick(tab tabWithTextBuffer, nick string) string {
	padamt := tab.Padlen(nick)
	if padamt > len(stripFmtChars(nick)) {
		return strings.Repeat(" ", padamt-len(stripFmtChars(nick))) + nick
	}
	return nick
}
