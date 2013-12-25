
package cirru

type coordObj struct {
  x, y int
}

// BufferObj represents a token like thing in Cirru
type BufferObj struct {
  Text string
  file *fileObj
  start, end coordObj
}

type tokenObj struct {
  class string
  buffer BufferObj
}

func tokenize(line []charObj) (tokens []tokenObj) {
  var buffer *BufferObj
  quoteMode := false
  escapeMode := false

  digestBuffer := func (asString bool) {
    if buffer != nil {
      
      newToken := tokenObj{}
      if asString {
        newToken.class = "string"
      } else {
        newToken.class = "text"
      }
      newToken.buffer = *buffer
      tokens = append(tokens, newToken)
      buffer = nil
    }
  }

  addBuffer := func (theChar charObj) {

    if buffer != nil {
      buffer.Text = buffer.Text + string(theChar.text)
      buffer.end.x = theChar.x
      buffer.end.y = theChar.y
    } else {
      start := coordObj{theChar.x, theChar.y}
      end := coordObj{theChar.x, theChar.y}
      text := string(theChar.text)
      file := theChar.file
      buffer = &BufferObj{text, file, start, end}
    }
  }

  for {
    if len(line) == 0 {
      break
    }
    char := line[0]
    line = line[1:]

    if quoteMode {
      if escapeMode {
        addBuffer(char)
        escapeMode = false
      } else {
        if char.isDoubleQuote() {
          digestBuffer(true)
          quoteMode = false
        } else if char.isBackslash() {
          escapeMode = true
        } else {
          addBuffer(char)
        }
      }
    } else {
      switch {
      case char.isBlank():
        digestBuffer(false)
      case char.isOpenParen():
        digestBuffer(false)
        newToken := tokenObj{}
        newToken.class = "openParen"
        tokens = append(tokens, newToken)
      case char.isCloseParen():
        digestBuffer(false)
        newToken := tokenObj{}
        newToken.class = "closeParen"
        tokens = append(tokens, newToken)
      case char.isDoubleQuote():
        digestBuffer(false)
        quoteMode = true
      default:
        addBuffer(char)
      }
    }
  }
  digestBuffer(false)
  return
}