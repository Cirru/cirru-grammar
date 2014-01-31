
package cirru

func parseText(line inline, args List) List {
  tokens := tokenize(line.line)

  var build func (byDollar bool) List
  build = func (byDollar bool) List {
    collection := List{}

    takeArgs := func () {
      if len(tokens) == 0 {
        if len(args) > 0 {
          for _, line := range args {
            if list, ok := line.(List); ok {
              dispersive := false
              if word, ok := list[0].(Token); ok {
                if word.Text == "," {
                  dispersive = true
                }
              }
              if dispersive {
                collection = append(collection, list[1:]...)
              } else {
                collection = append(collection, list)
              }
            }
          }
          args = List{}
        }
      }
    }

    takeArgs()

    for {
      if len(tokens) == 0 {
        if byDollar {
          if len(tokens) > 0 && tokens[0].class == "closeParen" {
            return collection
          }
        }
        break
      }
      cursor := tokens[0]
      tokens = tokens[1:]
      switch cursor.class {
      case "string":
        collection = append(collection, cursor.buffer)
      case "text":
        if cursor.buffer.Text == "$" {
          collection = append(collection, build(true))
        } else {
          collection = append(collection, cursor.buffer)
        }
      case "openParen":
        collection = append(collection, build(false))
      case "closeParen":
        return collection
      }
      takeArgs()
    }
    return collection
  }
  return build(false)
}