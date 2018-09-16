function! Buildgrid(expression, grid, index, lastLineIndex)
  let result =  matchstrpos(a:expression, "\\S", a:index)

  if result[1] == -1
    return
  else
    call add(a:grid, [result[0], result[1] - a:lastLineIndex])

    let newLineIndex = a:lastLineIndex

    if result[0] == "\n"
      let newLineIndex = result[1] + 1
    endif

    call Buildgrid(a:expression, a:grid, result[1]+1, newLineIndex)
  endif
endfunction

function! PasteGrid()
  let values = getreg()
  let grid = []
  call Buildgrid(values, grid, 0, 0)
  let lineCount = 0
  let pos = getcurpos()
  let currentRow = pos[1]
  let currentCol = pos[2]

  let currentOffset = pos[4]

  for note in grid
    if note[0] == "\n"
      let lineCount += 1
    else
      let row = lineCount + currentRow
      let col = note[1] + currentOffset
      call cursor(row, col)
      let command = "normal! r" . note[0]
      call execute(command)
    endif
  endfor
endfunction

nn <Leader>p :call PasteGrid()<CR>

function! GetOptionValue(name)
  call search(a:name)
  let lineNumber = getline('.')
  return split(lineNumber, "=")[1]
endfunction

let s:noteDictionary = {0: 'C-1', 1: 'C#-1', 2: 'D-1', 3: 'D#-1', 4: 'E-1', 5: 'F-1', 6: 'F#-1', 7: 'G-1', 8: 'G#-1', 9: 'A-1', 10: 'A#-1', 11: 'B-1', 12: 'C0', 13: 'C#0', 14: 'D0', 15: 'D#0', 16: 'E0', 17: 'F0', 18: 'F#0', 19: 'G0', 20: 'G#0', 21: 'A0', 22: 'A#0', 23: 'B0', 24: 'C1', 25: 'C#1', 26: 'D1', 27: 'D#1', 28: 'E1', 29: 'F1', 30: 'F#1', 31: 'G1', 32: 'G#1', 33: 'A1', 34: 'A#1', 35: 'B1', 36: 'C2', 37: 'C#2', 38: 'D2', 39: 'D#2', 40: 'E2', 41: 'F2', 42: 'F#2', 43: 'G2', 44: 'G#2', 45: 'A2', 46: 'A#2', 47: 'B2', 48: 'C3', 49: 'C#3', 50: 'D3', 51: 'D#3', 52: 'E3', 53: 'F3', 54: 'F#3', 55: 'G3', 56: 'G#3', 57: 'A3', 58: 'A#3', 59: 'B3', 60: 'C4', 61: 'C#4', 62: 'D4', 63: 'D#4', 64: 'E4', 65: 'F4', 66: 'F#4', 67: 'G4', 68: 'G#4', 69: 'A4', 70: 'A#4', 71: 'B4', 72: 'C5', 73: 'C#5', 74: 'D5', 75: 'D#5', 76: 'E5', 77: 'F5', 78: 'F#5', 79: 'G5', 80: 'G#5', 81: 'A5', 82: 'A#5', 83: 'B5', 84: 'C6', 85: 'C#6', 86: 'D6', 87: 'D#6', 88: 'E6', 89: 'F6', 90: 'F#6', 91: 'G6', 92: 'G#6', 93: 'A6', 94: 'A#6', 95: 'B6', 96: 'C7', 97: 'C#7', 98: 'D7', 99: 'D#7', 100: 'E7', 101: 'F7', 102: 'F#7', 103: 'G7', 104: 'G#7', 105: 'A7', 106: 'A#7', 107: 'B7', 108: 'C8', 109: 'C#8', 110: 'D8', 111: 'D#8', 112: 'E8', 113: 'F8', 114: 'F#8', 115: 'G8', 116: 'G#8', 117: 'A8', 118: 'A#8', 119: 'B8', 120: 'C9', 121: 'C#9', 122: 'D9', 123: 'D#9', 124: 'E9', 125: 'F9', 126: 'F#9', 127: 'G9'}

function! GetNotesFromConfiguration()
  let notesDefinition = GetOptionValue("Notes")
  let minMax = split(notesDefinition, "-")
  return range(minMax[0], minMax[1])
endfunction

function! GetNotes()
  let notes = []
  let definedNotes = GetNotesFromConfiguration()

  "TILME
  "let notes = sort(map(keys(s:noteDictionary), 'v:val + 0'), 'N')

  for noteNum in reverse(definedNotes)
    call add(notes, s:noteDictionary[noteNum])
  endfor

  return notes
endfunction

function! DrawKeyboard()
  hi Sharp   guibg=#111111
  hi Natural guibg=#222222
  let cursorpos = getcurpos()
  sign unplace *
  let notes = GetNotes()

  /PATTERN
  let currentLine = line(".") + 1

  let index = 0
  for noteName in notes
    let index += 1
    if match(noteName, "#") >= 0
      let defineCommand = "sign define " . noteName . " text=" . noteName[0:1] . " texthl=Sharp linehl=Sharp"
    else
      let defineCommand = "sign define " . noteName . " text=" . noteName[0:1] . " texthl=Natural linehl=Natural"
    endif
    call execute(defineCommand)
    let lineNumber = currentLine + index

    let placeCommand = "sign place " . index ." line=". lineNumber ." name=" . noteName . " buffer=" . bufnr('%')

    if line("$") < lineNumber
      call append(line("$"), "")
    endif

    call execute(placeCommand)
  endfor

  call cursor(cursorpos[1], cursorpos[2])
endfunction

hi Sharp   guibg=#111111
hi Natural guibg=#222222

function! LongestLine()
  let start = search("PATTERN") + 2
  let lines = getline(start, line("w$"))
  let longest = 0
  for line in lines
    if longest < strlen(line)
      let longest = strlen(line)
    endif
  endfor

  return longest
endfunction

function! DrawTimeline()
  let cursorpos = getcurpos()
  /PATTERN
  let timelineLine = line(".") + 1
  let division = GetOptionValue("Division")
  let length = GetOptionValue("Length")

  let longestLine = LongestLine()

  let totalLength = length * division

  if totalLength < longestLine
    let totalLength = longestLine
  endif

  let timeline = ""

  let i = 0
  while strlen(timeline) <= totalLength
    let i += 1
    let timelineAppendage = i
    let j = 0
    while strlen(timelineAppendage) <= (division - 1)
      let j += 1
      let timelineAppendage .= "-"
    endwhile

    let timeline .= timelineAppendage
  endwhile

  execute("setlocal colorcolumn=" . (totalLength + 1))

  call setline(timelineLine, timeline[:(totalLength - 1)])

  call cursor(cursorpos[1], cursorpos[2])
endfunction

"  Ö - evenly divided into 2 notes
"  ȯ - note on 0 - half of full length
"  o͘ - note on 1/2 - half of full length
"  ȍ̋ - divided into four - all four
"  ȍ - divided into four - first two
"  ő - divided into four - last two
"  ò - divide into four - first 
"  ó - divide into four - last 
"  o̎ - divide into four - middle two
"  o͆ - divide into four - first and last

let s:diacriticalMarks = [ "", "\u0308", "\u0307", "\u0358", "\u030B", "\u30F", "\u0300", "\u0301", "\u030E", "\u0346", "\u030D" ]

function! CycleDiacritical(amount)
  normal cl
  let currentChar = @"

  let i = 0
  let newValue = currentChar
  let baseChar = strpart(currentChar, 0, 1)

  let combineFn = 'baseChar . v:val'
  let baseMarks = map(copy(s:diacriticalMarks), combineFn)

  while len(baseMarks) > i
    let i += 1

    let index = i % len(baseMarks)
    let nextIndex = (index + a:amount) % len(baseMarks)
    let newValue = substitute(currentChar, baseMarks[index], baseMarks[nextIndex], "")

    if newValue != currentChar
      let i += 10000
    endif
  endwhile

  exec("normal a" . newValue)
endfunction

augroup beatspart
  autocmd!
  autocmd BufNewFile,BufRead *.part setlocal expandtab tabstop=2 shiftwidth=2 nowrap
  autocmd BufRead,BufWrite *.part call DrawKeyboard()
  autocmd BufRead,BufWrite *.part call DrawTimeline()
  autocmd BufNewFile,BufRead,BufEnter *.part setlocal virtualedit+=all
  autocmd BufLeave *.part setlocal virtualedit-=all
  autocmd BufNewFile,BufRead *.part nnoremap L :call CycleDiacritical(1)<CR>
  autocmd BufNewFile,BufRead *.part nnoremap H :call CycleDiacritical(-1)<CR>
  autocmd BufLeave *.part nunmap L
  autocmd BufLeave *.part nunmap H
augroup END

