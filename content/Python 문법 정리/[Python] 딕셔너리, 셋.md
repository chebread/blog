---
date: 2025-05-11
category: [Python 문법 정리]
published: true
fixed: false
---

## 딕셔너리
딕셔너리는 리스트와 비슷하다.

리스트와 다른 점은 항목의 순서를 따지지 않으며, 오프셋으로 항목을 선택할 수 없다.

**대신 값에 상응하는 고유한 키를 지정한다.**

**키는 모든 불변 타입(불리언, 정수, 부동소수점, 튜플, 문자열 등)이 될 수 있다.**

딕셔너리는 변경 가능하므로 요소를 추가, 삭제, 수정할 수 있다.

딕셔너리를 연관 배열(associative array), 해시(hash), 해시맵(hashmap), 딕트(dict)라고도 부른다.

## 딕셔너리 생성하기: {}
딕셔너리를 만드려면 중괄호 안에 콤마로 구분한 `키:값` 쌍을 지정한다.

`{}` 는 키:값 쌍이 없는 빈 딕셔너리다.

```python
dict = {
	"a": 1,
	"b": 2
}

a = {}
```

**딕셔너리의 변수 이름을 입력하면 모든 키와 값을 출력한다.**

```python
dict = {
	"a": 1,
	"b": 2
}
dict
{"a": 1, "b": 2}
```

리스트, 튜플, 딕셔너리의 마지막 항목에는 콤마를 입력하지 않아도 된다.

**파이썬에서 들여쓰기는 중요하지만, 키와 값을 입력할 때 중괄호 내에 반드시 들여쓰기를 할 필요는 없다.**

그러나 중괄호 내부에 들여쓰기를 하는 이유는, 가독성을 좋게 만들기 때문이다.

## 딕셔너리 생성하기: dict()
**키와 값을 dict() 함수 인수에 키워드 인자 형식으로 전달하여 딕셔너리를 생성할 수 있다.**

키워드 인자 형식으로 딕셔너리를 생성할 때는 키워드 인자는 오로지 문자열만 가능하기에 딕셔너리의 키는 문자열만 가능하다.

```python
acme_customer = dict(first="Wile", middle="E")
{"first": "Wile", "middle": "E"}
```

키워드 인자란,

```python
def greet(name, age):
    print(f"{name} is {age} years old")

greet(name="Alice", age=30)  # ← 여기가 키워드 인자
```
`name="Alice"`, `age=30`

이런 형식이 키워드 인자이다.

## 변환하기: dict()
dict() 함수를 사용해서 **두 값으로 이루어진 시퀀스**를 딕셔너리로 변환할 수 있다.

각 시퀀스의 첫 번째 항목은 키로, 두 번째 항목은 값으로 사용된다.

```python
lol = [['a', 'b'], ['c', 'd'], ['e', 'f']]
dict(lol)
{'a': 'b', 'c': 'd', 'e': 'f'}
```

두 항목의 리스트를 딕셔너리로 변환하는 예제이다.

```python
lot = [('a', 'b'), ('c', 'd'), ('e', 'f')]
dict(lot)
{'a': 'b', 'c': 'd', 'e': 'f'}
```

리스트, 튜플로 구성된 시퀀스를 딕셔너리로 변환할 수 있다.

```python
tol = (['a', 'b'], ['c', 'd'], ['e', 'f'])
dict(tol)
{'a': 'b', 'c': 'd', 'e': 'f'}
```

두 항목의 리스트로 구성된 튜플을 딕셔너리로 변환할 수 있다.

```python
tos = ('ab', 'cd', 'ef')
dict(tos)
{'a': 'b', 'c': 'd', 'e': 'f'}
```

두 문자열로 구성된 리스트를 딕셔너리로 변환할 수 있다.

```python
tos = ('ab', 'cd', 'ef')
dict(tos)
{'a': 'b', 'c': 'd', 'e': 'f'}
```

두 문자열로 구성된 튜플을 딕셔너리로 변환할 수 있다.

zip() 함수를 사용해서 두 항목의 시퀀스를 쉽게 생성할 수 있다.

## 항목 추가/변경: `[key]`
`dict[key] = value` 의 형태로 항목을 추가 또는 변경할 수 있다.

딕셔너리에 이미 존재하는 키라면 그 값은 새로운 값으로 대체된다.

키가 존재하지 않는다면 새 값과 키가 딕셔너리에 추가도니다.

리스트와 달리 딕셔너리를 할당할 때는 인덱스 범위가 벗어나는 예외에 대해 걱정할 필요가 없다.

단, 딕셔너리의 키를 변경할 수는 없으며, 기존에 존재하는 `키:값` 을 삭제하고 새로운 `키:값` 을 만들어야 한다.

```python
dic = {'name': 'pey', 'phone': '010-9999-1234', 'birth': '1118'}

dic['name'] = 'fey'
dic
{'name': 'fey', 'phone': '010-9999-1234', 'birth': '1118'}

dic['x'] = 3
dic
{'name': 'fey', 'phone': '010-9999-1234', 'birth': '1118', 'x': 3}
```

**딕셔너리의 키들은 반드시 고유해야 한다.**

**같은 키를 두 번 이상 사용하면 마지막 키:값으로 저장된다.**

```python
a = {
	'a': 1,
	'b': 2
	'a': 3
}
a
{ 'a': 3, 'b': 2 }
```

## 항목 얻기: `[key]`
`dict[key]` 로 키에 상응하는 값을 얻을 수 있다.

```python
dict = {
	'a': 2,
	'b': 3
}
dict['a']
2
```

만약 `dict[key]` 로 값을 얻을 때 딕셔너리에 해당 키가 없으면, 예외가 발생한다.

```python
dict = {
	'a': 2,
	'b': 3
}
dict['l']
Traceback (most recent call last):
  File "<python-input-50>", line 1, in <module>
    dict['l']
    ~~~~^^^^^
KeyError: 'l'
```

이 문제를 피하는 좋은 방법은, `in` 으로 해당 키가 존재하는지 확인하는 것이다.

```python
'l' in dict
False
```

## 항목 얻기: `get()`
`dict.get(키)` 함수를 사용하여 값을 얻을 수 있다.

```python
dict = {
	'a': 2,
	'b': 3
}
dict.get('a')
2
```

`dict.get(키, 옵션 값)` 처럼 옵션 값을 따로 지정하여 키가 존재하지 않을 때 옵션 값을 출력하게 할 수 있다.

```python
dict.get('l', 'foo')
'foo'
```

만약 옵션 값을 지정하지 않고 존재하지 않는 키의 값을 가져오면 `None`을 얻는다.

단, 대화식 인터프리터에서는 `None` 이 아니라 아무것도 출력하지 않는다.

```python
dict.get('l')
None
```

## 모든 키 얻기: keys()
딕셔너리의 모든 키를 가져오기 위해 `dict.keys()` 를 사용한다.

**keys() 는 이터러블 객체인 `dict_keys()` 를 반환한다.**

`dict_keys()` 이터러블 객체를 `list()` 변환 함수를 사용하여 리스트로 만들 수 있다.

```python
>>> a = {'name': 'pey', 'phone': '010-9999-1234', 'birth': '1118'}
>>> a.keys()
dict_keys(['name', 'phone', 'birth'])
>>> list(a.keys)
['name', 'phone', 'birth']
```

파이썬 3에서는 이터러블 값을 일반적인 리스트로 변환하기 위해 `list()` 함수를 사용한다.

`list()` 함수는 이터러블 객체를 리스트 타입으로 변환해주는 변환 함수이기 때문이다.

## 모든 값 얻기: values()
딕셔너리의 모든 값을 가져오기 위해 `dict.values()` 를 사용한다.

values() 는 이터러블 객체인 `dict_values()` 를 반환한다.

```python
>>> a = {'name': 'pey', 'phone': '010-9999-1234', 'birth': '1118'}
>>> a.values()
dict_values(['pey', '010-9999-1234', '1118'])
```

## 모든 키-값 얻기: items()
딕셔너리의 모든 키/값을 얻기 위해 `dic.items()` 을 사용한다.

items() 는 이터러블 객체인 `dict_items()` 를 반환한다.

각 키와 값은 튜플로 반환된다.

```python
>>> a.items()
dict_items([('name', 'pey'), ('phone', '010-9999-1234'), ('birth', '1118')])
>>> list(a.items())
[('name', 'pey'), ('phone', '010-9999-1234'), ('birth', '1118')]
```

## 길이 얻기: len()
`len()` 함수로 딕셔너리에 있는 키-값 쌍의 개수를 구할 수 있다.

```python
len(a)
3
```

빈 딕셔너리는 `0` 을 반환한다.

## 결합하기: `{**a, **b}`
**`**` 를 사용하여 딕셔너리를 두 개 이상 결합할 수 있다.**

여기서 `**` 는 제곱 연산자가 아니라 다른 용도로 사용된다.

`**` 연산자는 딕셔너리 `a`와 `b`를 푼 다음, 그 내용들을 새로운 딕셔너리에 넣는다.

결합하는 딕셔너리에 같은 키가 존재한다면 마지막 `키:값`이 저장된다.

이는 `같은 키를 두 번 이상 사용하면 마지막 키:값으로 저장된다.` 와 같은 원리이다.

```python
a = {'a': 1, 'b': 2}
b = {'b': 3, 'c': 4}
c = {**a, **b}
c
{'a': 1, 'b': 3, 'c': 4}
```

`**` 연산자를 사용하는 방법으로 두 개 이상 딕셔너리를 결합하게 되면,
얕은 복사를 수행하게 된다.

```python
a = {'x': [1, 2]}
b = {'y': 3}

c = {**a, **b}
c['x'].append(99)

print(a)
{'x': [1, 2, 99]} # a도 바뀜
```

딕셔너리 키/값 전체 사본을 깊은 복사하고 싶다면 `deepcopy()` 를 수행한다.

## 결합하기: update()
`update()` 함수는 한 딕셔너리의 키와 값들을 복사해서 다른 딕셔너리에 붙여준다.

`dict_a.update(dict_b)` 처럼 사용하며, `dict_b` 딕셔너리를 `dict_a` 에 삽입한다.

만약 삽입한 값 중에 기존 키:값이 있다면 마지막 `키:값`(두 번째 딕셔너리 값, 삽입한 딕셔너리 값)으로 변경된다.

```python
sample_dict = {
    '수학':80,
    '국어':90
}
 
## 수학 성적을 95로 수정, 영어 성적 70점 추가
sample_dict.update({'수학':95, '영어':70})
 
print(sample_dict)
{'수학':95, '국어': 90, '영어':70}
```

## 키와 del로 항목 삭제하기
`del dict[key]` 로 특정 `키:값` 항목을 삭제할 수 있다.

```python
>>> a = {1: 'a', 2: 'b', 'name': 'pey', 3: [1, 2, 3]}
>>> del a[1]
>>> a
{2: 'b', 'name': 'pey', 3: [1, 2, 3]}
```

## 키로 항목 가져온 뒤 삭제하기: pop()
`pop()` 은 `get()` 과 `del` 을 함께 사용하는 것과 같다.

딕셔너리에 있는 키와 `pop()` 의 인수가 일치한다면

해당 값을 반환한 뒤, 딕셔너리에서 해당 `키-값` 을 삭제한다.

만약 딕셔너리에 키가 존재하지 않으면 예외가 발생한다.

```python
```null
>>> sample_dict = {'a': 1, 'b': 2, 'c': 3, 'd': 4}
>>> sample_dict.pop('a')
1
>>> sample_dict
{'b': 2, 'c': 3, 'd': 4}
```

pop() 에 두 번째 인수를 지정하면 get() 과 같이 작동한다.

딕셔너리에 키가 존재하지 않으면 두 번째 인수가 출력된다.

```python
>>> sample_dict.pop('e', "딕셔너리에서 해당 key가 없습니다")
딕셔너리에서 해당 key가 없습니다
```

## 모든 항목 삭제하기: clear()
딕셔너리에 있는 키와 값을 모두 삭제하기 위해서는 `dict.clear()`를 사용하면 된다.

```python
a = { ... }
a.clear()
a
{}
```

또는, 빈 딕셔너리(`{}`)를 변수에 할당하면 모두 삭제할 수 있다.

```python
a = { ... }
a = {}
a
{}
```

## 키 맴버십 테스트: `in`
딕셔너리에 키가 존재하는지 알고 싶다면 `in` 을 사용한다.

**"값이 존재하는가"가 아님을 주의한다.**
**"키가 존재하는 가"를 묻는 것이다.****

`x in dict` 은 dict 딕셔너리에 `x` 라는 키가 존재하는가를 불리언 값으로 반환한다.

```python
a = { ... }
'a' in a
True
```

## 할당하기: `=`
리스트와 마찬가지로, 한 딕셔너리를 변수 두 곳에 할당했을 때, 한 딕셔너리를 변경하면 다른 딕셔너리도 같이 변경된다.

이를 얕은 복사라고 부른다.

```python
signals = { 'green': 'go', 'yellow': 'go faster', 'red': 'smile for the camera'}
save_signals = signals
signals['blue'] = 'confuse everyone'
save_singals
{ 'green': 'go', 'yellow': 'go faster', 'red': 'smile for the camera', 'blue': 'confuse everyone' }
```

## 복사: copy(), deepcopy()
copy() 를 사용하면 리스트와 마찬가지로 표면적으로는 깊은 복사인 것 처럼 동작하지만, 종속된 객체에 대해서는 얕은 복사가 수행된다.

```python
signals = { 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}
signals_copy = signals.copy()
signals
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}
signals_copy
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}

signals['red'][1] = 'sweat'
signals
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'sweat']}
signals_copy
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'sweat']} # 얕은 복사 발생
```

깊은 복사를 위해서는 deepcopy()를 사용해야 한다.

```python
import copy
signals = { 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}
signals_copy = copy.deepcopy(signals)
signals
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}
signals_copy
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']}

signals['red'][1] = 'sweat'
signals
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'sweat']}
signals_copy
{ 'green': 'go', 'yellow': 'go faster', 'red': ['stop', 'smile']} # 깊은 복사 발생
```

## 딕셔너리 비교
비교 연산자(`==, !=`)를 사용하여 딕셔너리를 비교할 수 있다.

단, 비교 연산자 `==, !=` 그 외의 비교 연산자는 사용할 수 없다.

파이썬은 키/값의 순서에 상관없이 딕셔너리의 키/값을 하나씩 비교한다.

```python
a = {1: 1, 2: 2, 3: 3}
b = {3: 3, 1: 1, 2: 2}

a == b
True

a <= b
Traceback...

a != b
False
```

## 순회하기: for in
**딕셔너리를 순회하면 키가 반환된다.**

즉, keys() 메서드 순회한 것과 같은 결과를 낸다.

```python
accusation = { 'a': 1, 'b': 2 }
for card in accusation:
	print(card)
a
b

for card in accusation.keys():
	print(card)
a
b
```

키가 아닌 값을 순회하려면 values() 메서드를 사용한다.

```python
for value in accusation.values():
	print(value)
1
2
```

키와 값을 모두 튜플로 반환하려면 items() 메서드를 사용한다.

```python
for item in accusation.items():
	print(item)
('a', 1)
('b', 2)
```

튜플 언패킹을 이용할 수도 있다.

```python
for card, contents in accusation.items():
	print(card, contents)
a, 1
b, 2
```

## 딕셔너리 컴프리헨션
리스트 컴프리헨션과 같이 딕셔너리 컴프리헨션이 있다.

```python
{키 표현식 : 값 표현식 for 표현식 in 이터러블}
```

```python
word = 'letters'
letter_counts = { letter: word.count(letter) for letter in word}
letter_counts
{'l': 1, 'e': 2, 't': 2, 'r': 1, 's': 1}

letter_counts = { letter: word.count(letter) for letter in set(word)}
{'l': 1, 'e': 2, 't': 2, 'r': 1, 's': 1}
```

리스트 컴프리헨션과 같이 딕셔너리 컴프리헨션 또한 if 문과 다중 for 문을 사용할 수 있다.

```python
{키 표현식: 값 표현식 for 표현식 in 이터러블 if 테스트}
```

```python
v = 'aeiou'
w = 'onomatopoeia'
v_counts = {letter: word.count(letter) for letter in set(word) if letter in v}
v_counts
{ ... }
```

## 셋(집합)
셋은 값은 버리고 키만 남은 딕셔너리와 같다.

셋의 각 요소는 유일해야 한다.

**셋에서 중복된 값은 자동으로 제거한다.**
**셋은 수학적 집합과 같은 개념이다.**

## 셋 생성하기
셋을 생성할 때는 `{}` 안에 콤마로 구분된 값을 넣는다.

```python
even_numbers = {0, 2, 3}
even_numbers
{0, 2, 3}
```

set() 함수로 빈 셋을 생성할 수 있다.

`{}` 는 빈 셋을 생성하지 않는다.

`{}` 는 빈 딕셔너리를 생성한다.

그래서 인터프리터는 빈 셋을 set()로 출력한다.

왜 그럴까? 파이썬에서 딕셔너리가 셋 보다 먼저 등장했기 때문이다.

```python
empty_set = set()
empty_set
set()
```

## 변환하기: set()
set() 함수를 사용하여 리스트, 문자열, 튜플, 딕셔너리 등 이터러블 객체의 중복된 값을 삭제하여 셋을 생성할 수 있다.

문자열을 set() 함수로 셋으로 만들어보자.

```python
set('letters')
{'l', 'e', 't', 'r', 's'}
```

'e'와 't'가 두 개씩 있어도 셋에는 하나만 저장된다.

셋에서 중복된 값은 자동으로 제거하기 때문이다.

리스트를 셋으로 만들어보자.

```python
set(['a', 'b'])
{'a', 'b'}
```

튜플을 셋으로 만들어보자.

```python
set(('a', 'b'))
{'a', 'b'}
```

딕셔너리에 set()을 사용하면 키만 셋으로 만들어진다.

```python
set({'a': 1, 'b': 2})
{'a', 'b'}
```

## 길이 얻기: len()
len() 함수는 셋의 원소의 개수를 반환한다.

```python
len(set({'a': 1, 'b': 2}))
2
```

빈 셋은 0개를 반환한다.

## 항목 추가하기: add()
add()를 활용하여 셋에 항목을 추가한다.

```python
s = set([1, 2, 3])
s.add(4)
{1, 2, 3, 4}
```

## 항목 삭제하기: remove()
`remove(값)` 함수는 값으로 셋의 항목 삭제한다.

```python
s = set([1, 2, 3])
s.remove(1)
{2, 3}
```

## 순회하기: for in
딕셔너리처럼 셋에 있는 모든 항목을 순회할 수 있다.

```python
s = set([1, 2, 3])
for a in s:
	print(a)
1
2
3
```

## 맴버십 테스트: in
`x in set` 으로 `set` 셋에 `x` 요소가 있는지에 대한 결과를 불리언 값으로 반환한다.

```python
s = set([1, 2, 3])
1 in s
True
```

## 교집합(intersection)
& 연산자와 intersection() 메서드를 사용해서 교집합을 구할 수 있다.

## 합집합(union)
| 연산자와 union() 메서드를 사용해서 합집합을 구할 수 있다.

## 차집합(difference)
`-` 연산자와 difference() 메서드를 사용해서 차집합을 구할 수 있다.

## 대칭 차집합(exclusive)
`^` 연산자나 symmertric_difference() 메서드를 사용해서 대칭 차집합을 구할 수 있다.

## 부분집합(subset)
<= 연산자나 Issubset() 메서드를 사용해서 부분집합을 구할 수 있다.

## 상위집합(superset)
`>=` 연산자나 issuperset() 메서드를 사용해서 상위집합을 구할 수 있다.

## 진부분집합(proper subset)
`<` 연산자를 사용해서 진부분집합을 구할 수 있다.

## 진상위집합(proper superset)
`>` 연산자를 사용해서 진상위집합을 구할 수 있다.

## 셋 컴프리헨션
셋 컴프리헨션은 리스트, 딕셔너리 컴프리헨션과 같은 모양이다.

```python
{표현식 for 표현식 in 이터러블}
```

또한 셋 컴프리헨션에 for 문과 if 문을 사용할 수 있다.

```python
{표현식 for 표현식 in 이터러블 if 테스트}
```

```python
a_set = {number for number in range(1, 6) if number & 3 == 1}
a_set
{1, 4}
```

## 불변 셋 생성하기: frozenset()
셋은 가변 객체이다.

불변 객체인 셋을 생성하려면 frozenset() 함수와 인수로 순회 가능한 객체(이터러블)을 사용한다.

```python
frozenset([3, 2, 1])
frozenset({1, 2, 3})
frozenset(set([2, 1, 3]))
fs = frozenset((1, 2, 3))
fs.add(4)
Traceback...
```

## 지금까지 배운 자료구조
지금까지 다음과 같은 자료구조를 배웠다.

- 대괄호를 사용한 리스트
- 콤마와 괄호를 사용한 튜플(괄호는 옵션이다)
- 중괄호를 이용한 딕셔너리 또는 셋

셋을 제외하고 모두 대괄호로 항목에 접근할 수 있다.

리스트와 튜플의 경우 대괄호에 들어가는 값이 정수 오프셋이고,
딕셔너리는 키다.

**셋은 오프셋과 키가 없다.**

## 자료구조 결합하기
다양한 자료구조를 결합해서 자료구조를 확장할 수 있다.

자료구조의 제한 사항은 데이터 타입 자체에 있다.

예를 들어 딕셔너리의 키는 불변흐기 때문에 리스트, 딕셔너리, 셋은 딕셔너리의 키가 될 수 없다.

그러나 튜플은 딕셔너리의 키가 될 수 있다.

예를 들면 어떤 장소를 GPS 좌표로 인덱싱할 수 있다.

```python
houses = {
(44, -93, 285): 'My house'
}
```