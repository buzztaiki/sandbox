# https://narimiran.github.io/nim-basics/

proc doArray() =
 var
   a: array[3, int] = [5, 7, 9]
   b = [5, 7, 9]
   # c = []  # error
   d: array[7, string]

 const m = 3
 let n = 5

 var e: array[m, char]
 # var f: array[n, char] # error


proc doSequence() =
 var
   e1: seq[int] = @[]
   f = @["abc", "def"]
 var
   e = newSeq[int]()


proc doSlice() =
 let j = ['a', 'b', 'c', 'd', 'e']
 echo j[1]
 echo j[^1]                      # 最後の要素
 echo j[0 .. 3]                  # スライス (3を含む)
 echo j[0 ..< 3]                 # スライス (3を含まない)

proc doTuple() =
  let n = ("Banana", 2, 'c')
  echo n

  var o = (name: "Banana", weight: 2, rating: 'c')
  o[1] = 7
  o.name = "Apple"
  echo o


proc doFunctionCall() =
  proc plus(x, y: int): int =
    return x + y

  proc multi(x, y: int): int =
    return x * y

  let
    a = 2
    b = 3
    c = 4

  # 第一引数をレシーバみたいに使える。おもしろ。
  echo a.plus(b) == plus(a, b)
  echo c.multi(a) == multi(c, a)


  echo a.plus(b).multi(c)
  echo c.multi(b).plus(a)

import strutils
proc doUseModule() =
  let
    a = "My string with whitespace."
    b = '!'
  echo a.split()
  echo split(a)
  echo a.toUpperAscii()
  echo b.repeat(5)

doArray()
doSequence()
doSlice()
doTuple()
doFunctionCall()
doUseModule()
