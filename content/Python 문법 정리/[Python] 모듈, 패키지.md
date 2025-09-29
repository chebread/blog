---
date: 2025-06-03
category: [Python 문법 정리]
published: true
fixed: false
---

## 모듈과 import 문
모듈이란 파이썬 코드 파일이다.
개발자는 특별한 작업 없이 코드를 모듈로 사용할 수 있다.

## 모듈 임포트하기
import 문을 사용하여 다른 모듈의 코드를 참조할 수 있다.
import 문을 사용하여 간단하게 모듈을 임포트 할 수 있다.
확장자 .py를 제외한 파이썬 파일의 이름을 입력한다.

```python
# fast.py
from random import choice
places = ["a", "b", "c", "d", "e"]
def pick():
	return choice(places)
```

```python
# lunch.py
import fast
place = fast.pick()
print(place)
```

lunch.py를 실행하면 메인 프로그램은 fast 모듈에 접근해서 fast() 함수를 실행한다.
메인 프로그램이란, 파이썬 인터프리터가 가장 먼저 실행하도록 직접 지정된 스크립트 파일이다.

import 문 이후, fast.py에 있는 모든 객체는 이름 앞에 `fast.`를 붙여 메인 프로그램에서 사용할 수 있다.
**즉, import로 불러온 모듈에 내포된 모든 객체는 `모듈 이름.모든 객체` 이런 방식으로 import를 불러온 파일에서 모든 불러온 모듈에 대한 객체를 사용할 수 있다는 말이다.**

```python
# fast2.py
places = ['a', 'b', 'c', 'd', 'e']
def pick():
	import random
	return random.choice(places)
```
이렇게 함수 내부에서 import 문을 사용할 수도 있다.
이렇게 함수 내부에서 import 문을 사용하는 경우, 코드 사용이 그 함수 내부에서만 제한된다.
그러므로, 모든 import 문은 파일 코드의 맨 위에 두는 것을 파이썬 개발자들은 선호한다.

**더불어, 모듈을 임포트 하면 모듈을 한 번 실행하고, 이를 sys.module에 추가한다.**

`import this` 하면 파이썬의 철학이 출력된다.
모듈만 임포트 하면 어떻게 문자열이 출력될까?
this.py 모듈에는 print()가 포함되어 있다.
모듈을 임포트 하면 모듈을 한 번 실행하고, 이를 sys.module에 추가한다.
즉, 임포트 하면 한 번 그 모듈 파일을 실행하니까, print()가 실행되니까 출력되는 것이다.

## 다른 이름으로 모듈 임포트하기
```python
# fast3.py
import fast as f
place = f.pick()
print(place)
```
이렇게 에일리어스(alias)를 사용할 수 있다.
`import 모듈 이름 as 축약할 이름` 이런 문법으로 축약할 수 있다.
단, 지정한 축약 이름(에일리어스, alias)으로만 해당 모듈을 참조할 수 있다. 원래 모듈 이름인 `fast`로는 직접 참조할 수 없다.

## 필요한 모듈만 임포트하기
모듈 전체나 필요한 부분만 임포트할 수 있다.
`from 모듈 이름 import 필요한 객체` 이렇게 해서 필요한 객체만 모듈에서 가져올 수 있다.
```python
# fast4.py
from fast import pick
place = pick()
print(place)
```

```python
# fast5.py
from fast import pick as who_cares
place = who_cares()
print(place)
```

## 패키지
파이썬 애플리케이션을 좀 더 확장하기 위해 모듈을 패키지(package)라는 파일과 모듈 계층 구조에 구성할 수 있다.
패키지는 .py 파일을 포함한 하위 디렉터리이다.
단, 패키지 디렉터리 안에 `.py` 파일 외의 다른 종류의 파일(예: 데이터 파일, 텍스트 파일, 이미지 파일, 컴파일된 확장 모듈 등)이 포함되어 있어도 패키지로 취급되는 데는 전혀 문제가 없다.
- `__init__.py` 파일이 있다면 해당 디렉터리는 명시적으로 "일반 패키지"로 취급된다.
- `__init__.py` 파일이 없고 `.py` 파일들만 있다면, Python 3.3 이상에서는 "네임스페이스 패키지"의 일부로 간주된다.

`from 패키지 이름 import 모듈 이름` 이렇게 패키지 내부의 모듈에 접근할 수 있다.
패키지 내부에 패키지가 있는 경우는 `from 패키지 이름.패키지 이름 import 모듈 이름` 이렇게 모듈에 접근할 수 있다.

```python
# file tree
choice/
	fast.py
	advice.py
```

```python
# choice/fast.py
from random import choice

place = ['a', 'b', 'c', 'd', 'e']
def pick():
	return choice(place)
```

```python
# choice/advice.py
from random import choice
answers = ['a', 'b', 'c']
def give()"
	return choice(answers)
```

```python
# question.py
from choice import fast, advice

print(fast.pick())
print(advice.give())
```
현재 디렉터리에서 choice라는 패키지(디렉터리)를 찾는다. 그리고 그 안에 있는 fast.py와 advice.py 모듈(파일)을 찾는다.

> 근데 `from 모듈 이름 import 필요한 객체` 이렇게 해서 필요한 객체만 모듈에서 가져올 수 있는데, 패키지와 어떻게 구분하나요? 패키지도 from 패키지 이름 import 모듈 이름 이렇게 가져오는데?
> => 그냥 파이썬이 알아서 구분합니다.

파이썬 3.3 이전 버전인 경우, 코드 하위 디렉터리에 `__init__.py` 파일이 하나 더 필요하다(하위 디렉터리를 패키지로 만들어준다.) `__init__.py` 파일은 빈 파일이다. 그러나, 파이썬 3.3 이전 버전에서는 이 파일을 포함하는 디렉터리만 패키지로 간주하기 때문에 필요하다.
단, 파이썬 3.3 이후 버전은 자동으로 패키지라는 것을 간주한다.

## 모듈 탐색 경로
파이썬 인터프리터가 보는(임포트하는) 모든 위치를 보려면 표준 sys 모듈을 임포트해서 path 리스트를 살펴본다.
```python
import sys
for place in sys.path
	print(place)

...
```

출력 첫 줄의 공백은 현재 디렉터리를 뜻하는 빈 문자열이다.

중복된 이름의 모듈이 있다면 첫 번째로 검색된 모듈을 사용한다.
**만약, 우리가 random 이라는 모듈을 정의하고, 이 모듈이 표준 라이브러리 random을 찾기 전에 검색 경로에 이미 있다면 표준 라이브러리 random 모듈을 사용할 수 없다.**

코드 내에서 탐색 경로를 수정할 수 있다.
**파이썬이 다른 것보다 먼저 `/my/modules` 디렉터리에서 탐색하길 원한다고 가정하면 다음과 같이 코드를 추가한다.**
```python
import sys
sys.path.insert(0, "my/modules")
```

## 상대/절대 경로 임포트
파이썬은 상대(absolute)/절대(relative) 경로 임포트를 지원한다.

**"절대 임포트"는 모듈의 이름을 항상 `sys.path`의 시작점부터 전체적으로 지정하는 방식이고, "상대 임포트"는 현재 모듈의 위치를 기준으로 경로를 지정하는 방식이라고 이해하면 된다.**
- **절대 임포트 구문**: 점(`.`)으로 시작하지 않는 경로 (예: `import package.module`)
- **상대 임포트 구문**: 점(`.`)으로 시작하는 경로 (예: `from . import module`)

주의할 점은, 절대 임포트와 상대 임포트 모두 상대 경로를 사용한다는 거다.
파이썬 `import`는 모듈 이름 또는 점(`.`)으로 연결된 패키지 경로를 사용한다.

상대 임포트는 현재 모듈의 위치를 기준으로 다른 모듈이나 패키지를 가져오는 방식이다.
"절대 임포트(absolute import)"는 파이썬에서 모듈이나 패키지를 가져올 때 사용되는 가장 기본적인 방식으로, 가져오려는 대상의 전체 경로를 `sys.path`에 있는 최상위 디렉터리 중 하나로부터 시작하여 명시하는 방법이다.

- `import rougarou` => 이것은 절대 임포트(absolute import) 방식이다.
- `from . import rougarou` => 이것은 상대 임포트 방식이다. 여기서 점(`.`)은 "현재 패키지(current package)"를 의미한다.

- rougarou.py 파일이 메인 프로그램을 실행한 파일과 같은 디렉터리에 있는 경우, from . import rougarou을 사용하여 상대 경로 임포트를 할 수 있다.
- 상위 디렉터리에 있는 경우, from .. import rougarou을 사용한다.
- 상위 디렉터리의 creatures라는 디렉터리에 있는 경우, from ..creatures import rougarou를 사용한다.
- (현재 디렉터리)와 ..(부모 디렉터리) 표기법에 대한 것은 유닉스에서 사용했다.

## 네임스페이스 패키지
```
north
	critters
		wendigo.py

south
	critters
		rougarou.py
```

이런 상황에서는, `from critters import wendigo, rougarou` 이렇게 단일 디렉터리 패키지를 공동으로 사용하는 것처럼 모듈을 가져올 수 있다.

네임스페이스 패키지는 하나의 패키지가 **물리적으로 여러 개의 서로 다른 디렉터리에 나뉘어 존재**할 수 있도록 하는 기능이다. 마치 여러 조각들이 모여 하나의 전체 그림을 이루는 것처럼, 서로 다른 위치에 있는 디렉터리들이 동일한 이름의 패키지 공간을 공유하게 된다.

## 모듈 vs 객체
언제 코드를 모듈로 사용해야 하는가? 아니면 객체로 사용해야 하는가?
모듈과 객체는 여러면에서 비슷하게 보인다.
(모듈 이름.객체 이렇게 사용하고 객체도 객체 이름.속성or메서드 이렇게 사용하기 때문이다.)

한 프로그램에서 임포트한 모듈 사본은 하나만 있다.
이를 이용하여 임포트한 모든 코드에 전역값을 저장할 수 있다.

=>
"한 프로그램에서 임포트한 모듈 사본은 하나만 있다."
이 말은 정확합니다. 파이썬 프로그램이 실행될 때, 특정 모듈은 최초에 한 번만 로드(load)됩니다.

1. 어떤 모듈이 처음으로 `import` 되면, 파이썬은 해당 모듈 파일을 찾아서 실행하고, 그 결과로 생성된 모듈 객체를 만듭니다.
2. 이 모듈 객체는 `sys.modules`라는 특별한 딕셔너리(일종의 캐시)에 모듈 이름을 키로 하여 저장됩니다. 예를 들어 `import math`를 하면 `sys.modules['math']`에 `math` 모듈 객체가 저장됩니다.
3. 이후 프로그램의 다른 부분에서 같은 모듈을 다시 `import` 하려고 하면, 파이썬은 파일을 다시 읽고 실행하는 대신 `sys.modules`에 이미 저장된 모듈 객체를 그대로 가져와 사용합니다.

따라서, 한 프로그램 내에서는 특정 모듈에 대해 단 하나의 모듈 객체만이 존재하게 됩니다. 모든 `import` 문은 이 동일한 객체를 가리키게 됩니다.

"이를 이용하여 임포트한 모든 코드에 전역값을 저장할 수 있다."
이 말도 정확합니다. 위에서 설명한 것처럼 모든 곳에서 동일한 모듈 객체를 참조하기 때문에, 그 모듈 객체의 속성(attribute) 값을 변경하면, 해당 모듈을 임포트한 다른 모든 코드에서도 변경된 값을 보게 됩니다. 모듈의 속성은 해당 모듈을 사용하는 모든 코드에 대해 일종의 "전역 변수"처럼 작동할 수 있습니다.

"`import math; math.pi = 3.0` 이렇게 하면 `math.pi`라는 값이 다 바뀐다는 것인가?"
**네, 그렇습니다! 현재 실행 중인 파이썬 프로그램 내에서는 `math.pi`의 값이 3.0으로 변경되며, 이후 `math.pi`를 참조하는 모든 곳에서 이 변경된 값을 사용하게 됩니다.**
**그러나, 사용자님이 예시로 드신 `math.pi = 3.0`처럼 표준 라이브러리나 잘 알려진 외부 라이브러리의 값을 이렇게 직접 수정하는 것은 매우 위험하며 일반적으로 강력히 권장되지 않는 방식입니다.**

## 파이썬 표준 라이브러리
파이썬은 배터리 포함(batteries included)이라는 모토가 있다.
표준 라이브러리 모듈은 핵심 언어가 늘어나는 것을 피하기 위해 분리되어 있다.

## 누락된 키 처리하기: setdefault()와 defaultdict()
**`setdefault()` 함수는 `get()` 함수와 같지만, 키가 누락된 경우 딕셔너리에 항목을 할당할 수 있다.**
존재하는 키에 다른 기본값을 할당하려면 키에 대한 원래 값이 반환되고 아무것도 바뀌지 않는다.

```python
periodic_table = { 'Hydrogen': 1, 'Helium': 2 }
periodic_table
{ 'Hydrogen': 1, 'Helium': 2 }

carbon = periodic_table.setdefault('Carbon', 12)
carbon
12

periodic_table
{ 'Hydrogen': 1, 'Helium': 2, 'Carbon': 12 }

helium = periodic_table.setdefault('Helium', 947)
helium
2
periodic_table
{ 'Hydrogen': 1, 'Helium': 2, 'Carbon': 12 }
```

defaultdict() 함수도 비슷하다. 다른 점은 딕셔너리를 생성할 때 모든 새 키에 대한 기본값을 먼저 지정한다는 것이다. 이 함수의 인수는 함수다.

## 항목 세기: Counter()
표준 라이브러리에는 항목을 셀 수 있는 여러 함수가 있다.
```python
from collections import Counter
breakfast = ['spam', 'spam', 'eggs', 'spam']
breakfast_counter = Counter(breakfast)
breakfast_counter
Counter({ 'spam': 3, 'eggs': 1 }) # `collections.Counter`는 파이썬의 내장 `dict` (딕셔너리) 클래스의 서브클래스(자식 클래스)이다.
```

most_common 함수는 모든 요소를 내림차순으로 반환한다. 혹은 숫자를 입력하는 경우, 그 숫자만큼의 상위 요소를 반환한다.
```python
breakfast_counter.most_common()
[('spam', 3), ('eggs', 1)]

breakfast_counter.most_common(1)
[('spam', 3)]
```

카운터를 결합-빼기-교집합-합집합 할 수 있다.
```python
lunch = ['eggs', 'eggs', 'bacon']
lunch_counter = Counter(lunch)

breakfast_counter + lunch_counter
Counter({'spam': 3, 'eggs': 3, 'bacon': 1})

breakfast_counter - lunch_counter
Counter({'spam': 3})

lunch_counter - breakfast_counter
Counter({'bacon': 1, 'eggs': 1})

breakfast_counter & lunch_counter
Counter({'eggs': 1})

breakfast_counter | lunch_counter
Counter({'spam: 3, 'eggs': '2', 'bacon': 1})
```

## 키 정렬하기: OrderedDict()
파이썬 3.7 버전 부터는 딕셔너리의 키는 추가된 순서대로 그 키 순서를 유지한다.
OrderedDict는 파이선 3.7 이전 버전에서 키 순서를 유지하려 할 때 유용하다.
```python
from collections import OrderedDict
quotes = OrderedDict([
	('a', 2),
	('b', 3)
])

for stooge in quotes:
	print(stooge)
'a'
'b'
```

## 스택 + 큐 == 데크
데크(deque)는 스택과 큐의 기능을 모두 가진 출입구가 양 끝에 있는 큐다.
데크는 시퀀스의 양 끝으로부터 항목을 추가하거나 삭제할 때 유용하게 쓰인다.
여기에서 회문(palindrome)(앞에서부터 읽으나 뒤에서부터 읽으나 같은 구문) 인지 확인하기 위해 양쪽 끝에서 중간까지 문자를 확인한다.
popleft() 함수는 데크로부 터 왼쪽 끝의 항목을 제거한 후, 그 항목을 반환한다.
pop() 함수는 오른쪽 끝의 항목을 제거한 후, 그 항목을 반환한다.
양쪽 끝에서부터 이 두 함수가 중간 지점을 향해서 동작한다.
양쪽 문자가 서로 일치한다면 단어 중간에 도달할 때까지 데크를 팝(pop)한다.

```python
def palindrome(word):
	from collections import deque
	dq = deque(word)
	while len(dq) > 1:
		if dq.popleft() != dq.pop():
			return False
	return True

palindrome('a')
True

palindrome('racecar)
True

palindrome('')
True

palindrome('radar')
True

palindrome('halibut')
False
```

데크의 간단한 예제를 살펴봤다.
회문 코드를 더 간단하게 작성하고 싶다면, 한 문자열을 반전해서 비교하면 된다.
파이썬은 문자열에 대한 reverse() 메서드가 없지만, 다음과 같이 슬라이스로 문자열을 반전할 수 있다.
```python
def another_palindrome(word):
	return word == word[::-1]

another_palindrome('radar')
True

another_palindrome('halibut')
False
```

## 코드 구조 순회하기: itertools
itertools는 특수 목적의 이터레이터 함수를 포함하고 있다.
for ... in 반복문에서 이터레이터 함수를 호출할 때 함수는 한 번에 한 항목을 반환하고 호출 상태를 기억한다.

chain() 함수는 순회가능한 인수들을 차례로 반복한다.
```python
import itertools
for item in itertools.chain([1, 2], ['a', 'b'])
	print(item)
1
2
a
b
```

cycle() 함수는 인수를 순환하는 무한 이터레이터다.
```python
import itertools
for item in itertools.cycle([1, 2])
1
2
1
2
...
```

accumlate() 함수는 축적된 값을 계산한다. 기본으로 합계를 계산한다.
```python
import itertools
for item in itertools.accumulate([1, 2, 3, 4]):
	print(item)
1
3
6
10
```

accumulate() 함수의 두 번째 인수를 전달하여, 합계를 구하는 대신 이 함수를 사용할 수 있다.
이 함수는 두 개의 인수를 취하여 하나의 결과를 반환한다.
이 예제는 축적된 곱을 계산한다.
```python
import itertools
def multiply(a, b):
	return a * b
	
for item in itertools.accumlate([1, 2, 3, 4], multiply):
	print(item)
1
2
6
24
```

itertools 모듈은 많은 함수를 제공하며, 조합 및 순열을 위한 함수도 있다.
직접 작성하는 구현 시간을 절약할 수 있다.

## 깔끔하게 출력하기: pprint()
지금까지 모든 예제를 print() 함수를 사용해서 출력하거나 대화식 인터프리터에서 변수 이름으로 출력했다.
출력된 결과를 읽기 힘든 경우 pprint() 함수와 같은 멋진 프린트가 필요하다.

print() 함수의 출력 결과는 다음과 같다.
pprint() 함수는 가독성을 위해 요소들을 정렬하여 출력한다.
```python
from pprint import pprint
quotes = Ordered([
	('a', 1),
	('b', 2)
])

print(quotes)
Ordered(['a', 1], ['b', 2])

pprint(quotes)
{'a': 1, 'b': 2}
```

## 임의값 얻기: random
random.choice() 는 주어진 것(리스트, 튜플, 딕셔너리, 문자열)의 인수에서 값을 임의로 반환하다.

한 번에 둘 이상의 값을 얻으려면 sample() 함수를 사용한다.

어떤 범위에서 임의의 정수를 얻으려면 choice() 와 sample()을 range()와 같이 사용하거나 randint()와 randrange()를 사용한다.

randrange()는 range()와 같은 인수를 가진다. 시작(포함), 끝(제외), 스탭(옵션값)

random()을 사용하여 0.0과 0.1 사이의 임의의 실수를 얻는다.

## 배터리 장착: 다른 파이썬 코드 가져오기
외부 패키지를 사용할 수 있다.
