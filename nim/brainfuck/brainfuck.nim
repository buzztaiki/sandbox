import os
import system/io

type OutOfRangeError = object of ValueError
type State = object
  data: array[30000, char]
  index: int

proc next(state: var State, n: int) =
  state.index += n
proc value(state: State): char =
  state.data[state.index]
proc setValue(state: var State, value: char) =
  state.data[state.index] = value
proc incr(state: var State, n: int) =
  state.setValue(char(ord(state.value) + n))

proc findCloseBracket(code: string, start: int): int =
  var n = 0
  for i in start ..< code.len:
    if code[i] == '[':
      n += 1
    elif code[i] == ']':
      if n == 0:
        return i
      n -= 1
  raise newException(OutOfRangeError, "")

proc evalCode(code: string, state: var State): void

proc evalLoop(code: string, state: var State) =
  while ord(state.value()) != 0:
    evalCode(code, state)
    
proc evalCode(code: string, state: var State) =
  var i = 0
  while i < code.len:
    case code[i]:
      of '>': state.next(1)
      of '<': state.next(-1)
      of '+': state.incr(1)
      of '-': state.incr(-1)
      of '.': stdout.write(state.value())
      of ',': state.setValue(stdin.readChar())
      of '[':
        let j = findCloseBracket(code, i+1)
        evalLoop(code[i+1 .. j-1], state)
        i = j
      else:
        discard
    i += 1

proc main() =
  var code = paramStr(1)
  var state: State
  evalCode(code, state)

# e.g. brainfuck "+++++++++[>++++++++>+++++++++++>+++++<<<-]>.>++.+++++++..+++.>-.------------.<++++++++.--------.+++.------.--------.>+."
main()
