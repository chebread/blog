---
date: 2025-05-06
category: [Python 문법 정리]
published: true
fixed: false
---

## 반복문
코드를 한 번 이상 실행하려고 할 때 반복문(loop)가 필요하다.
파이썬은 두 개의 반복문 while과 for이 있다.

## while문
`while 조건: ...` 형태로 쓰며, 조건이 참이면 while문을 실행하고, 조건이 거짓이면 while문이 종료된다.
while 문이 종료되면 그 다음 줄이 실행된다.

```python
count = 1
while count <= 5:
	print(count)
	count += 1 # count++ 같이 증감 연산자는 파이썬에서 없다.
```

## while문 break
while문에서 break를 사용하면 while문을 중단하여 탈출할 수 있다.

```python
while True:
	stuff = input()
	if stuff == 'q':
		break
	print(stuff.capitalize())
```

## while문 continue
다음 반복으로 즉시 건너뛸 때 continue를 사용한다.

```python
while True:
	value = input()
	if value == 'q':
		break
	number = int(value)
	if number % 2 == 0:
		# 짝수
		continue
	print(number, "squared is", number * number)
```

## while-else문
while 문을 작성하고, 그 아래에 else 문을 작성할 수 있다.
반복문이 break를 쓰지 않고 자연스럽게 종료됬다면 else문이 실행된다.
반복문이 break에 의해 종료되었다면 else문이 실행되지 않는다.

```python
>>> countdown = 5
>>> while countdown > 0:
...     print(countdown)
...     countdown -= 1
...     if input() == '중단':
...         break
... else:
...     print('발사!')
... 
5
4
3
중단
```

else문은 break 문에 의해 반복문이 중단되지 않고 모든 항목을 순회했는지 확인할 때 유용하다.

while 문을 무언가를 찾는 것으로 생각하고, 찾지 못했다면 else문이 호출된다고 생각하면 좀 더 이해하기 쉬울 것이다.

## for (in) 문
파이썬에서 이터레이터(iterator)는 유용하게 자주 쓰인다.
이터레이터는 자료구조를 순회할 수 있도록 해준다.
이터레이터는 데이터가 메모리에 맞지 않더라도 데이터 스트림을 처리할 수 있도록 허용해준다. => 터레이터가 모든 데이터를 한 번에 메모리에 올리지 않고, 필요할 때마다(요청할 때마다) 하나씩 값을 생성해서 반환한다는 뜻

일단, 문자열을 순회하는 방법만 살펴보자.

for (in) 문으로 문자열을 쉽게 순회할 수 있다.

for (in) 문은 `for 변수 in 리스트(또는 튜플, 문자열): ...` 형태이다.

for in 문이라고 부르는 이유는, 항상 for 문은 in을 사용하기 때문이다.

그러나 for 문에서 활용되는 in은 맴버십 연산자가 아니다.

for문에서의 in은 단순히 멤버십 연산자처럼 True/False를 반환하는 게 아니다.

**`for x in a:` 라고 작성하게 되면, "a의 모든 요소를 하나씩 x에 대입하며 반복한다"는 특별한 의미로 작동하게 된다.**

**문법적으로는 같은 in 이지만, for 문에서는 반복을 위한 키워드로 쓰이는 것이다.**

`in` 좌변에 오는 변수에 `in` 우변에 오는 변수의 각 요소가 전달되어 for in 반복문 내부에서 해당 변수가 활용된다.

```python
word = "abcd"
for letter in word:
	print(letter)
a
b
c
d
```

이것은 while 문으로는 이렇게 나타낸다.

```python
word = "abcd"
offset = 0
while offset < len(word):
	print(word[offset])
	offset += 1
a
b
c
d
```

while 문 보다는 for in 문의 방법이 더 파이써닉한 방법이다.

## for 문 break
break 문을 사용하여 for 문을 중단할 수 있다.

## for 문 continue
continue 문을 사용하여 for 문을 다음 반복으로 즉시 건너뛸 수 있다.

## for-else 문
for 문을 작성하고, 그 아래에 else 문을 작성할 수 있다.
반복문이 break를 쓰지 않고 자연스럽게 종료됬다면 else문이 실행된다.
반복문이 break에 의해 종료되었다면 else문이 실행되지 않는다.

```python
>>> for x in [1, 2, 3, 4]:
...     print(x)
... else:
...     print("리스트의 원소를 모두 출력했어요")
... 
1
2
3
4
리스트의 원소를 모두 출력했어요
```

else문은 break 문에 의해 반복문이 중단되지 않고 모든 항목을 순회했는지 확인할 때 유용하다.

for 문을 무언가를 찾는 것으로 생각하고, 찾지 못했다면 else문이 호출된다고 생각하면 좀 더 이해하기 쉬울 것이다.

## range() 함수
`range(start, stop, step)` 의 형태로 숫자 시퀀스를 생성할 수 있다.

range() 함수는 순회 가능한(이터러블) 객체를 반환한다.

> range() 함수의 반환값이 이터러블이라서, list(range(...)) 하면 리스트가 생성되는 거다.

start, start + step, start + step + step, ... 처럼 step 만큼 start에 더해져서 stop과 더해진 값이 같아지면 중단하고, start 부터 stop과 같아지기 직전의 step이 더해진 값 까지를 숫자 시퀀스로 반환한다.

> range() 함수를 "특정 구간의 숫자의 범위를 만들어주는 함수"라고 이해하는 것보다, "start 부터 step을 계속 더해가서 stop 까지 간 시퀀스 값을 반환하는 함수"라고 이해하는 것이 낫다.

start == stop 이면 빈 시퀀스 반환한다.
start > stop이고 step > 0 이면 빈 시퀀스 반환한다.
step = 0 이면 예외가 발생한다.
start를 생략시 0으로 취급된다.
step을 생략시 1으로 취급된다.

**start > 0 이면 start ~ stop - 1 까지, step은 1 이상**
**start <= 0 이면 start ~ stop + 1 까지, step은 1 이하**

range(2, -1, 1) 이면 2, 2 + 1, 2 + 1 + 1, 2 + 1 + 1 + 1, ... 이렇게 계속해서 step 만큼 더해져도 stop 까지 가지를 못한다.
그래서 start > stop 이면 step은 음수여야 한다는 것이다.
range(2, -1, -1) 이면 2, 2 - 1, 2 - 1 - 1, ... 이렇게 계속해서 step 만큼 더해지면 stop(-1) 까지 갈 수 있다.

step이 꼭 1 or -1일 필요는 없다.
step이 2면 2씩 더해지는 거고, -2면 -2씩 더해지는 거다.