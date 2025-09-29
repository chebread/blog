---
date: 2025-06-20
category: [Python 문법 정리]
published: true
fixed: false
---

## 파일과 디렉터리
파일(file)은 일련의 바이트이며 일부 파일 시스템(file system)에 저장되어 있는 파일 이름(filename)으로 접근한다.
디렉터리(directory)는 파일과 다른 디렉터리의 모음집이다.
폴더는 컴퓨터가 GUI를 얻었을 때 생겼으며, 디렉터리와 동의어다.

## 파일 입출력
데이터를 가장 간단하게 지속하려면 보통 파일(plain file)을 사용한다.
이것을 플랫 파일(flat file)이라고 부르기도 한다.
파일은 단지 파일 이름으로 저장된 바이트 시퀀스다.
파일에서 데이터를 읽어서 메모리에 적재하고, 메모리에서 파일로 데이터를 쓴다.
파이썬은 이러한 작업을 쉽게 만들어 준다.
이러한 파일 연산은 익숙하고 인기 있는 유닉스 같은 운영체제를 모델로 만들어졌다.

'plain file'은 데이터베이스 파일이나 워드 프로세서 문서(`.docx`)처럼 복잡한 내부 계층, 레코드, 필드, 메타데이터 구조를 갖지 않는다. 파일 자체가 자신의 구조를 정의하지 않으며, 단지 데이터(바이트)가 처음부터 끝까지 순차적으로 나열되어 있을 뿐이다.

Python은 이러한 'plain file'을 다룰 때, 그 안의 바이트를 어떻게 해석할지에 따라 크게 두 가지 종류(텍스트 파일 (Text File), 바이너리 파일 (Binary File))로 구분합니다. 이 두 종류 모두 'plain file'의 범주에 속합니다.

텍스트 파일: 파일 안의 바이트들이 특정 인코딩(예: UTF-8) 규칙에 따라 인간이 읽을 수 있는 문자로 해석될 수 있는 플레인 파일이다.
바이너리 파일: 파일 안의 바이트들이 문자뿐만 아니라 이미지, 소리, 실행 코드 등 기계가 해석해야 하는 원시(raw) 데이터를 포함하는 플레인 파일이다.

## 생성하기/열기: open()
open() 함수는
- 기존 파일 읽기
- 새 파일 쓰기
- 기존 파일에 추가하기
- 기존 파일 덮어쓰기
를 수행한다.

```python
fileobj = open(filename, mode)
```
이렇게 호출하며,
- fileobj: open() 함수에 의해 반환되는 파일 객체다.
- filename: 파일의 문자열 이름이다.
- mode: 파일 타입과 파일로 무엇을 할지 명시하는 문자열이다.

mode의 첫 번째 글자는 작업(operation)을 명시한다.
- r: 파일 읽기
- w: 파일 쓰기
	- 파일이 존재하지 않으면 파일을 생성하고, 파일이 존재하면 파일을 초기화(기존 내용을 모두 삭제)한 후에 재작성한다.
- x: 파일 쓰기
	- 파일이 존재하지 않을 경우에만 해당한다.
- a: 파일 추가하기
	- 파일이 존재하면 파일의 끝에서부터 쓴다.

mode의 두 번째 글자는 파일 타입을 명시한다.
- t or 아무것도 명시하지 않음: 텍스트 파일 타입
- b: 이진 파일 타입

> 파일 확장자(extension)의 핵심 역할은 해당 파일의 종류(format)가 무엇인지 알려주는 역할을 한다. 단, 확장자는 단지 파일 이름의 일부일 뿐인 관례(convention)이자 약속이다. 확장자가 파일의 실제 내용을 보장하지는 않는다. 예를 들어, `music.mp3` 파일의 이름을 `music.txt`로 바꿀 수 있다.
> 
> 확장자가 없는 파일도 똑같은 파일이다. 확장자가 없다는 것은 그저 이름에 '꼬리표'가 붙어있지 않다는 의미일 뿐, 파일의 본질(바이트 시퀀스)이 달라지는 것은 아니다.

## 파일 닫기: close()
파일을 열은 후 다 사용했다면, 사용한 메모리를 해제(free)하기 위해 파일을 닫아야 한다.
파일은 `fileobj` 객체의 `close()` 메서드로 해제한다.
```python
fout = open('oops.txt', 'wt')
fout.close()
```

## 텍스트 파일 쓰기: print()
print() 함수에서 file=fileobj 인수를 추가하여 내용을 쓸 수 있다.
print() 함수에서 file 인수가 없다면, print() 함수는 터미널인 표준 출력에 내용을 쓴다.
`print()` 함수에서 `file` 인수를 생략하면, file 인수에는 기본값인 `sys.stdout` 으로 설정된다.
`sys.stdout` 이란`sys` 모듈에 정의된 객체로, 표준 출력(Standard Output)을 의미한다. '표준 출력'은 프로그램이 실행될 때 약속된 기본 출력 통로이다. 우리가 터미널이나 명령 프롬프트에서 파이썬 코드를 실행하면, 이 '표준 출력' 통로는 기본적으로 터미널 화면에 연결된다.
```python
fout = open('oops.txt', 'wt')
print('Oops, I created a file.', file=fout)
fout.close()
```
단, print 함수로 파일에 내용 써도 `\n` 추가된다.

print() 함수의 인자중 `sep` 은 여러 인수를 print할 때, 그 인수 사이에 구분자를 지정하는 것이다. 기본 값은 ` ` 공백이다.
`sep='\n'` 이렇게 하면 여러 인수를 print할 때 인수 사이에 개행을 추가하여 print 할 수 있다.
```python
fout = open('oops.txt', 'wt')
print('a', 'b', 'c', sep='\n', end="", file=fout)
fout.close()

# oops.txt
a
b
c
```

- [ ] print 함수, input 함수는 따로 정리하기

## 텍스트 파일 쓰기: write()
write() 함수를 사용하여 파일에 내용을 쓸 수 있다.
write() 함수는 `fileobj.write(쓸 내용)` 이렇게 사용한다.
write() 함수는 텍스트 모드에서 파일에 쓴 문자(character)의 수를 반환한다.
write() 함수는 이진 모드에서 파일에 쓴 바이트(byte)의 수를 반환한다.
```python
poem = '''...
...
...
'''
fout = open('relativity', 'wt')
fout.write(poem)
```

print() 를 write() 처럼 작동하려면 print() 에 다음 두 인수를 전달한다.
- sep=""
- end=""
```python
poem = '''...
...
...
'''
fout = open('relativity, 'wt')
print(poem, file=fout, sep='', end='')
```

파일에 쓸 문자열이 크면, 특정 단위(chunk)로 나눠 파일에 쓴다.
```python
poem = '''...
...
...
'''
fout = open('relativity', 'wt')
size = len(poem)
offset = 0
chunk = 100
while True:
	if offset > size:
		break
	fout.write(poem[offset:offset+chunk])
	offset += chunk
fout.close()
```

relativity 파일이 있는데 이것을 모르고 실수로 덮어쓰고 싶지 않다면 x 모드를 사용한다.
```python
fout = open('relativity', 'xt')
Traceback...
```

이를 다음과 같이 예외로 처리할 수 있다.
```python
try:
	fout = open('relativity', 'xt')
	fout.write('stomp stomp stomp')
except FileExistsError:
	print('exists!')

'exists!'
```

## 텍스트 파일 읽기: read()
read() 함수를 인수 없이 호출하여 한 번에 파일 전체 내용을 읽을 수 있다.
그러나, 대형 파일로 이 작업을 수행할 때 많은 메모리가 소비되므로 주의해야 한다.
read() 함수를 사용하기 위해서는, `'r'` 모드로 파일을 열여야 한다.
쓰기 전용 모드인 `'w'`, `'a'`, `'x'`로 파일을 열면 `.read()`를 호출할 수 없다.
단, 파일 모드에 `+`를 추가하면 업데이트(읽고 쓰기 모두)를 위해 파일을 연다.
따라서 `+`가 포함된 모드에서는 `.read()`를 사용할 수 있다.
```python
fin = open('relativity', 'rt')
poem = fin.read()
fin.close()
```

read() 함수에서 한 번에 얼마만큼 읽을 수 있는지 크기를 제한할 수 있다.
read() 함수에 최대 읽을 문자 수를 인수로 입력한다.
```python
poem = ''
fin = open('relativity', 'rt')
chunk = 100
while True:
	fragment = fin.read(chunk)
	if not fragment:
		break
	poem += fragment
fin.close()
```

read() 함수는 파일의 끝에 도달하여 더 이상 읽을 바이트가 없을 때 `.read()` 메서드는 빈 문자열(`''`)을 반환한다.

`while` 루프가 계속 돌다가 파일 내용을 모두 읽고 나면, `fragment = fin.read(chunk)` 라인에서 `fragment` 변수에는 `''` 값이 저장된다.
파이썬에서는 모든 객체를 참 또는 거짓으로 평가할 수 있다.
not fragment 를 하게 되면 fragment가 공백이면 True가 된다. 파이썬에서 거짓으로 평가되는 것 외에는 모두 참으로 평가된다. 그러므로 공백이 아니면 not fragment는 False가 될 것이다.

## 텍스트 파일 읽기: readline()
readline() 함수를 사용하여 파일을 줄 단위로 읽을 수 있다.
```python
poem = ''
fin = open('relativity', 'rt')
while True:
	line = fin.readline()
	if not line:
		break
	poem += line
fin.close()
```
텍스트 파일의 빈 줄의 길이는 1이고('\n'), 이것을 True로 인식한다.
주의할 점, 텍스트 파일의 빈 줄은 시각적으로 비어있을 뿐, 실제로는 한 줄을 구성하는 `\n` 문자를 포함하고 있다.
readline() 함수는 read() 함수처럼 파일 끝에 도달했을 때, False로 간주하는 빈 문자열을 반환한다.

## 텍스트 파일 읽기: 이터레이터
텍스트 파일을 가장 읽기 쉬운 방법은 이터레이터를 사용하는 것이다.
`open()` 함수는 파일 객체를 반환한다.
이 파일 객체를 순환(iterate)하면, 그 파일의 내용을 한 줄씩 순서대로 읽어올 수 있다.
```python
poem = ''
fin = open('relativity', 'rt')
for line in fin:
	poem += line
fin.close()
```

## 텍스트 파일 읽기: readlines()
readlines() 함수는 한 번에 모든 줄을 읽고, 한 줄로 된 문자열가 담긴 리스트(list)를 반환한다.
```python
fin = open('relativity', 'rt')
lines = fin.readlines()
fin.close()
for line in lines:
	print(line, end="")
...
...
```
단, 각 줄에는 개행문자가 포함되어 있을 수도 있고 없을 수도 있다.
`readlines()` 함수는 단지 파일에 있는 내용을 정직하게 그대로 읽어오는 것만 수행한다.

## 이진 파일 쓰기: write()
mode를 `b` 로 변경시, 파일은 이진 모드로 열린다.
그러면 문자열 대신 바이트를 읽고-쓸 수 있다.
```python
bdata = bytes(range(0, 256))
fout = open('bfile', 'wb')
fout.write(bdata)
fout.close()
```

텍스트 파일 처럼 특정 단위 청크로 이진 데이터를 쓸 수 있다.
```python
bdata = bytes(range(0, 256))
fout = open('bfile', 'wb')
size = len(bdata)
offset = 0
chuck = 100
while True:
	if offset > size:
		break
	fout.write(bdata[offset:offset+chuck])
	offset += chuck
fout.close()
```

## 이진 파일 읽기: read()
파일을 rb 모드 또는 `+` 모드로 열어서 이진 파일을 읽을 수 있다.
```python
fin = open('bfile', 'rb')
bdata = fin.read()
fin.close()
```

## 자동으로 파일 닫기: with
열려 있는 파일을 닫지 않았을 때, 파이썬은 이 파일이 더 이상 참조되지 않는 것을 확인한 뒤 파일을 닫는다.
이것은 함수 안에 파일을 열어놓고 이를 명시적으로 닫지 않더라도 함수가 끝날 때 자동으로 파일이 닫힌다는 것을 의미한다.
그러나 오랫동안 작동하는 함수 혹은 메인 프로그램에 파일을 열어 놓았다면, 파일에 쓰는 것을 마치기 위해 명시적으로 파일을 닫아야 한다.
파이썬은 파일을 여는 것과 같은 일을 수행하는 컨텍스트 매니저(context manager)가 있다.
파일을 열때 `with 표현식 as 변수` 형식을 사용한다.
`as`는 표현식을 as 뒤에 선언한 특정 변수에 담아 `with` 블록 안에서 사용할 수 있게 해주는 역할을 한다.
```python
poem = ''
with open('relativity', 'wt') as fout:
	fout.write(poem)
```
컨텍스트 매너저 코드 블록의 코드 한 줄이 실행되고 나서(잘 수행되거나 문제가 있으면 예외 발생) 자동으로 파일을 닫아준다.

`with open(...) as fout:` 구문은 컨텍스트 매니저(Context Manager) 라고 불리는 매우 강력한 기능이다.
`with` 문을 쓰는 진짜 이유는 코드가 완벽하게 실행되지 않았을 때, 즉 중간에 에러(예외, Exception)가 발생했을 때 드러난다.
```python
print("\n--- with 문을 사용하는 경우 ---")

try:
    with open('error_test_with.txt', 'w') as f:
        print("with 블록 진입, 파일이 열렸습니다.")
        # 똑같은 에러 발생!
        num_written = f.write(12345)
        print("이 문장은 절대 출력되지 않습니다.")

except TypeError as e:
    print(f"TypeError가 발생했습니다: {e}")

f = open('error_test_with.txt', 'r')
print(f"에러 발생 후, 파일은 정상적으로 생성되었고 내용은 비어있습니다. 파일 상태: {f.closed}")
f.close()
print("파일을 다시 열었다가 닫았습니다.")
```
`f.write(12345)`에서 `TypeError`가 발생하는 순간, 프로그램의 실행 흐름이 `except` 블록으로 점프해 버버린다.
그 아래에 있던 `f.close()` 코드는 아예 실행될 기회조차 얻지 못했다.
그 결과 파일은 계속 열려있는 상태로 방치된다.
이것은 심각한 문제를 일으킬 수 있다.
그러므로, `with` 문은 단순히 `close()`를 안 쓰기 위한 편의 기능이 아니라, 어떤 얘기치 못한 상황에서도 파일을 안전하게 닫는 것을 보장하는 매우 중요한 안전장치이다.
단, `with` 문은 파일만을 위한 문법이 아니며, `with` 문은 파이썬의 컨텍스트 관리 구문(Context Management Statement) 이라는 범용적인 문법이다.

더불어, 주의할 점은 파이썬의 `with` 문은 새로운 스코프(scope)를 만들지 않는다는 것이다.
따라서 `with` 블록 안에서 생성하고 값을 할당한 변수(`today_string`)는 `with` 블록이 끝난 후에도 해당 스코프(예: 함수 내부, 전역 등)에서 계속해서 접근하고 사용할 수 있다.

## 파일 위치 찾기: tell(), seek()
tell() 함수는 파일 시작 위치에서 현재 오프셋을 바이트 단위로 반환한다.
seek() 함수는 다른 바이트 오프셋으로 위치를 이동한다.

## 메모리 메핑: mmap 모듈
파일을 읽고 쓰는 것의 대안은 표준 mmap 모듈로 파일을 메모리에 메핑하는 것이다.
`mmap`모듈은 운영체제의 가상 메모리 기능을 활용하여, 디스크 상의 파일을 메모리의 배열처럼 직접 다루게 해주는 강력하고 효율적인 고급 파일 처리 기법이다.

## 파일 명령어
파이썬은 다른 언어들처럼 유닉스의 파일 연산 패턴을 지니고 있다.
파이썬이 os.path 모듈과 새로운 pathlib 모듈로 이러한 작업을 처리하는 방법을 살펴본다.

## 존재 여부 확인하기: exists()
파일 혹은 디렉터리가 실존하는지 확인하기 위해 os.path 모듈의 exists() 함수를 사용한다.
exists() 함수의 반환값은 불리언이다.
exists() 함수의 인수로 상대 경로와 절대 경로를 사용할 수 있다.
- **절대 경로**: `C:\project\data\sales.csv`
- **상대 경로**: `data\sales.csv` 또는 `.\data\sales.csv`
	- `.`:  현재 디렉터리를 나타낸다.
	- `..`: 상위(부모) 디렉터리를 나타낸다.
여기서 말하는 절대/상대 경로는 우리가 통상적으로 사용하는 절대/상대 경로를 의미한다.
즉, 여기서 말하는 절대/상대 경로는 파이썬의 특수한 절대/상대 임포트에서의 절대/상대 경로가 아니다.
```python
import os
os.path.exists('oops.txt')
True
os.path.exsits('./oops.txt')
True
os.path.exists('.')
True
os.path.exists('..')
True
```

## 파일 확인하기: isfile()
os.path 모듈의 isfile() 함수는 인수로 절대/상대 경로를 전달하여 해당 경로의 것이 실존하는지, 파일인지 확인한다.
isfile() 함수는 실존하지 않거나 파일이 아니라면, False를 반환한다.
isfile() 함수는 실존하거나 파일이라면, True를 반환한다.
```python
import os
name = 'oops.txt'
os.path.isfile(name)
True
```

## 디렉터리 확인하기: isdir()
os.path 모듈의  isdir() 함수는 인수로 절대/상대 경로를 전달하여 해당 경로의 것이 실존하는지, 디렉터리인지 확인한다.
isdir() 함수는 실존하지 않거나 디렉터리가 아니라면, False를 반환한다.
isdir() 함수는 실존하거나 디렉터리라면, True를 반환한다.
```python
import os
name = 'oops.txt'
os.path.isdir(name)
False

os.path.isdir('.')
True
```

## 절대 경로 확인하기: isabs()
os.path 모듈의 isabs() 함수는 인수가 절대 경로인지 확인한다.
```python
import os
name = 'oops.txt'
os.path.isabs(name)
False

os.path.isabs('/big/fake')
True
```

## 파일 복사하기: copy()
shutil 모듈의 copy(file1, file2) 함수는 file1 파일을 file2에 복사한다.
단, 디렉터리는 복사할 수 없다.
```python
import shutil
shutil.copy('oops.txt', 'ohno.txt')
```

## 파일 이름 바꾸기: rename()
os 모듈의 rename() 함수는 파일 이름을 변경한다.
```python
import os
os.rename('ohno.txt', 'ohwell.txt')
```

## 연결하기: link(), symlink()
link() 함수는 하드 링크를 생성한다.
symlink() 함수는 심볼릭 링크를 생성한다.

## 권한 바꾸기: chmod()
os 모듈의 chmod() 는 파일 권한(permission)을 변경한다.
파일 권한은 사용자에 대한 읽기, 쓰기, 실행권한, 사용자가 속한 그룹과 나머지에 대한 권한이 있다.
이 명령은 사용자, 그룹, 나머지 권한을 묶어서 압축된 8진수의 값을 취한다.

```python
import os
os.chmod('oops.txt', 0o400)
```
oops.txt 파일을 이 파일의 소유자(파일을 생성한 사용자)만 읽을 수 있도록 만들었다.

이러한 8진수 값을 사용하기보다 심벌을 사용하고 싶다면, stat 모듈을 임포트하여 이렇게 쓰면 된다.
```python
import os
import stat
os.chmod('oops.txt', stat.S_IRUSR)
```

## 소유권 바꾸기: chown()
os 모듈의 chown() 함수는 숫자로된 사용자 아이디(uid)와 그룹 아이디(gid)를 지정하여 파일의 소유자와 그룹에 대한 소유권(ownership)을 바꿀 수 있다.
```python
uid = 5
gid = 22
os.chown('oops.txt', uid, gid)
```

## 파일 지우기: remove()
os 모듈의 remove() 함수를 사용하여 파일을 삭제할 수 있다.
삭제하려는 파일이 존재하지 않을 경우, `os.remove()` 함수는  예외(Exception)를 발생시킨다.
```python
import os
os.remove('oops.txt')
```

## 디렉터리 명령어
대부분 운영체제에서 파일은 디렉터리(폴더)의 계층 구조 안에 존재한다.
이러한 모든 파일과 디렉터리의 컨테이너는 파일 시스템(볼륨(volume))이다.
표준 os 모듈은 이러한 운영체제의 특성을 처리하고, 조작할 수 있는 함수를 제공한다.

## 디텍터리 생성하기: mkdir()
os 모듈의 mkdir() 함수를 사용하여 디렉터리를 생성한다.
```python
import os
os.mkdir('poems')
```

## 디렉터리 삭제하기: rmdir()
os 모듈의 rmdir() 함수를 사용하여 디렉터리를 삭제한다
```python
import os
os.rmdir('poems')
```

## 콘텐츠 나열하기: listdir()
os 모듈의 listdir() 함수를 사용하여 디렉터리의 내용을 리스트 타입으로 나열한다.
`os.listdir()`는 오직 해당 디렉터리 바로 아래에 있는 파일과 하위 디렉터리들의 이름까지만 나열한다.
하위 디렉터리 안으로 더 들어가서 그 안에 있는 파일이나 폴더명까지 보여주지는 않는다.
현재부터 하위의 모든 내용을 나열하기 위해서는 `os.walk()` 함수를 사용한다.
```python
import os
os.listdir('poems')
```

## 현재 디렉터리 위치 바꾸기: chdir()
os 모듈의 chdir() 함수는 현재 디렉터리의 위치를 다른 디렉터리의 위치로 바꾼다.
`os.chdir(새로운_경로)` 형태로 함수를 호출하며, 호출시 스크립트의 `현재 위치` 가 `새로운_경로`로 즉시 변경된다.
```python
import os
os.chdir('poems')
```

## 일치하는 파일 나열하기: glob()
glob 모듈의 glob() 함수는 복잡한 정규 표현식이 아닌, 유닉스 셸 규칙을 사용하여 일치하는 파일이나 디렉터리 이름을 검색해준다.
규칙은 다음과 같다.
- 모든 것에 일치: `*` (re 모듈의 `*` 와 같다.)
- 한 문자에 일치: `?`
- `a` or `b` or `c` 문자에 일치: `[abc]`
- `a` or `b` or `c` 를 제외한 문자에 일치: `[!abc]`

`m` 으로 시작하는 모든 파일이나 디렉터리를 찾는다.
```python
import glob
glob.glob('m*')
...
```

두 글자로 된 파일이나 디렉터리를 찾는다.
```python
glob.glob('??')
...
```

m으로 시작하고 e로 끝나는 여덟 글자의 단어를 찾는다.
```python
glob.glob('m??????e')
...
```

k or l or m으로 시작하고 e로 끝나는 단어를 찾는다.
```python
glob.glob('[klm]*e')
...
```

## 경로 이름
모든 컴퓨터는 계층적 파일 시스템을 사용한다.
디렉터리는 파일과 디렉터리를 다양한 형태로 포함한다.
특정 파일 혹은 디렉터리를 참조하려면 경로 이름(pathname)이 필요하다.
최상위 루트(root) 혹은 현재 디렉터리는 어떤 경로에 도달하기 위한 기준 위치다.
보통 슬래시(`/`)와 백슬리시(`\`)를 많이 혼동한다.
쉽게 기억하는 방법은 슬래시는 앞으로 기울어지고, 백슬래시는 뒤로 기울어진다.
유닉스와 macOS, 웹 주소는 경로 구분 기호로 슬래시를 사용한다.
윈도우는 백슬래시를 사용한다.

> 빌게이츠는 IBM이 첫 번째 PC를 요청했을 때, 'MS-DOS'를 갖기 위해 5만 달러에 QDOS라는 운영체제를 구입했다.
> 여기에서 명령행 인수에 슬래시를 사용하는 CP/M이라는 것을 흉내 냈다.
> MS-DOS가 나중에 폴더를 추가할 때 백슬래시를 사용해야 했다.

파이썬은 슬래시를 경로 구분 기호로 사용한다.
윈도우에서 백슬래시를 사용할 수 있지만, 파이썬에서 백슬래시는 이스케이프 문자다.
그러므로 윈도우에서는 백슬래시를 두 번 입력하거나, 원시(raw) 문자열(`r''`)을 사용해야 한다.
```python
win_file = 'eek\\urk.txt'
win_file = r'eek\urk2.txt'
```

경로 이름을 작성하는 방법은 다음과 같다.
- 경로 구분 문자(슬래시 또는 백슬래시)를 사용한다.
- os.path.join() 으로 경로 이름을 만든다.
- pathlib 모듈을 사용한다.

## 절대 경로 얻기: abspath()
os.path 모듈의 abspath() 함수는 상대 경로를 절대 경로로 만들어 준다.

현재 디렉터리가 `/usr/gaberlunzie` 고, `oops.txt` 파일이 거기에 있다면, 다음과 같이 절대 경로를 만들 수 있다.
```python
import os
os.path.abspath('oops.txt')
'/usr/gaberlunzie/oops.txt'
```

## 심볼릭 링크 얻기: realpath()
os.path 모듈의 realpath() 함수는 주어진 경로에 포함된 모든 심볼릭 링크(symbolic link) 요소를 해석하여, 해당 파일이나 디렉터리의 실제 물리적인 절대 경로를 반환한다.

## 경로 이름 작성하기: os.path.join()
os.path.join() 을 호출하여 운영체제에 적합한 경로 구분 문자로 경로를 결합할 수 있다.
`os.path.join`은 파일 시스템 경로를 다룰 때 발생할 수 있는 오류를 줄여주고, 코드의 호환성을 높여주는 필수적인 함수다.
```python
import os
win_file = os.path.join('eek', 'urk')
win_file2 = os.path.join(win_file, 'snort.txt')
```
이것을 macos나 linux에서 실행 시,
```python
'eek/urk/snort.txt'
```
이것을 windows에서 실행 시,
```python
'eek\\urk\\snort.txt'
```

그러나, 이처럼 같은 코드가 실행되는 위치에 따라 다른 결과가 나온다면 문제가 될 수 있다.
이러한 문제를 위해서 pathlib 모듈을 사용할 수 있다.

## pathlib 모듈
pathlib 모듈은 os.path 모듈의 대안으로 사용할 수 있다.
파일 시스템 경로 이름을 문자열로 취급하는 대신 Path 객체를 도입하여 파일 시스템 경로 이름을 더 높은 수준으로 처리한다.
Path 클래스를 사용하여 Path 객체를 만든 후, 경로를 '/' 문자가 아닌 그냥 / 슬래시로 묶는다.
파일 이름이나 경로 이름 문자열처럼 file_path를 open() 함수의 인수로 사용할 수 있다.
```python
from pathlib import Path
file_path = Path('eek') / 'urk' / 'snort.txt'

file_path
PosixPath('eek/urk/snort.txt')

print(file_path)
'eek/urk/snort.txt'

file_path.name
snort.txt

file_path.suffix
'.txt'

file_path.stem
'snort'
```

## BytesIO와 StringIO
(pass)