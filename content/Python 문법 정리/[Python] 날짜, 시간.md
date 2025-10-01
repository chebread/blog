---
date: 2025-06-20
category: [Python 문법 정리]
published: true
fixed: false
---

## 날짜, 시간 파이썬 표준 라이브러리
개발자는 날짜와 시간에 많은 노력을 기울인다.
이번 2장에서는 프로그래머들이 부딪히는 몇 가지 문제에 대해 살펴보고, 이를 좀 더 단순하게 만드는 트릭과 모범 사례에 대해 살펴본다.
파이썬 표준 라이브러리는 datetime, time, calendar, dateutil 등 시간과 날짜에 대한 여러 가지 모듈이 있다.
일부 중복되는 기능이 있어서 조금 헷갈린다.

## 윤년
윤년(leap year)은 특정한 시간 주기다.
윤년은 4년에 한 번씩 온다.
100년 마다 오는 해는 윤년이 아니지만, 400년 마다 오는 해는 윤년이라는 사실을 알고 있는가?
```python
import calendar
calendar.isleap(1900)
False

calendar.isleap(1996)
True

calendar.isleap(1999)
False

calendar.isleap(2000)
True

calendar.isleap(2002)
False

calendar.isleap(2004)
True
```

- 1년은 365.242196일이다.
	- 지구가 태양 주위를 한 번 돈 후, 시작된 지점의 축에서 약 1/4 회전한다.
- 4년마다 하루씩 추가한다.
	- 약 365.242196 - 0.25 = 364.992196일이다.
- 100년마다 하루씩 뺀다.
	- 약 364.992196 + 0.01 = 365.002196일이다.
- 400년마다 하루씩 추가한다.
	- 약 365.002196 - 0.0025 = 364.999696
거의 365일에 가깝다. 윤초(leap seconds)라는 개념도 있다.

## datetime 모듈
표준 datetime 모듈은 날짜와 시간을 처리한다
datetime 모듈은 여러 메서드를 가진 4개의 클래스를 제공한다.
- date: 년, 월, 일
- time: 시, 분, 초, 마이크로초
- datetime 날짜, 시간
- timedelta: 날짜 및 시간 간격
## date 클래스
datetime 모듈의 date 클래스는 day, month, year를 인자로 받으며, datetime.date 객체(date 객체)를 반환한다.
`datetime.date(...)` 객체는 여러 속성을 제공한다.
```python
from datetime import date
halloween = date(2019, 10, 31) # 객체
halloween
datetime.date(2019, 10, 31)

halloween.day
31

halloween.month
10

halloween.year
2019
```

## date 객체.isoformat() 메서드
date 객체는 isoformat() 메서드를 사용하여 `yyyy-mm-dd` 형식의 날짜를 출력할 수 있다.
```python
from datetime import date
halloween = date(2019, 10, 31) # 객체
halloween.isoformat()
2019-10-31
```

## date 클래스.today() 메서드
date.today() 메서드는 오늘 날짜를 출력한다.
```python
from datetime import date
date.today()
datetime.date(2019, 4, 5)
```

## date 객체.weekday() 메서드
datetime.date 객체의 weekday() 메서드를 사용하면 해당 날짜에 대한 요일을 가져올 수 있다.
weekday() 메서드에서는 0이 월요일이고, 6이 일요일이다.

## date 객체.isoweekday() 메서드
isoweekday() 메서드에서는 1이 월요일이고, 7이 일요일이다.

## timedelta 객체
datetime 모듈의 timedelta 객체를 사용하여 date 객체에 시간 간격을 더할 수 있다.
날짜의 범위는 date.min 부터 date.max 까지다.
결과적으로 천문학적인 날짜는 계산할 수 없다.
```python
from datetime import date
from datetime import timedelta

now = date.today()
one_day = timedelta(days=1)
tomorrow = now + one_day
tomorrow
datetime.date(2019, 4, 6)

now + 17 * one_day
datetime.date(2019, 4, 22)

yesterday = now - one_day
yesterday
datetime.date(2019, 4, 4)
```

## time 객체
datetime 모듈의 time 객체는 하루의 시간을 나타내는 데 사용된다.
time 클래스는 인자로 hour, minute, second, microsecond 를 받고, `datetime.time` 객체(time 객체)를 반환한다.
`datetime.time(...)` 객체는 여러 속성을 제공한다.
만약, time 클래스의 인자를 입력하지 않으면 초기 인수는 0으로 간주된다.
주의할 점은, 컴퓨터는 ms를 정확히 계산할 수 없다.
ms 측정의 정확성은 hw와 os의 많은 요소에 따라 달라진다.
```python
from datetime import time
noon = time(12, 0, 0) # ms는 0으로 간주
noon
datetime.time(12, 0)

noon.hour
12

noon.minute
0

noon.second
0

noon.microsecond
0
```

## datetime 클래스
datetime 모듈의 datetime 객체는 날짜와 시간을 모두 포함한다.
datetime 클래스는 인수로 year, month, day, hour, minute, second, microsecond 를 받으며,  datetime.datetime 객체(datetime 객체)를 반환한다.
```python
from datetime import datetime
some_day = datetime(2019, 1, 2, 3, 4, 5, 6)
some_day
datetime.datetime(2019, ..., 6)
```

## datetime 객체.isoformat() 메서드
datetime 객체에도 isoformat() 메서드가 있다.
isoformat() 메서드의 결과값의 `T` 는 날짜와 시간을 구분한다.
```python
from datetime import datetime
some_day = datetime(2019, 1, 2, 3, 4, 5, 6)
some_day.isoformat()
'2019-01-02T03:04:05.0000006'
```

## datetime 클래스.now() 메서드
datetime 클래스의 now() 메서드로 현재 날짜와 시간을 얻을 수 있다.
```python
from datetime import datetime
now = datetime.now()
now
datetime.datetime(2019, ...)

now.year
2019

...

now.microsecond
580562
```

## datetime 클래스.combine() 메서드
datetime 클래스의 combine(date 객체, time 객체) 메서드로 date 객체와 time 객체를 datetime 객체로 병합할 수 있다.
```python
from datetime import datetime, time, date # 한꺼번에 모듈을 불러올 수 있다
noon = time(12)
this_day = date.today()
noon_today = datetime.combine(this_day, noon)
noon_today
datetime.datetime(...)
```

## time 모듈
datetime 모듈의 time 객체와 별도의 time 모듈이 존재한다.
헷갈림을 주의해야 한다.

**절대 시간을 나타내는 한 가지 방법은 어떤 시작점 이후 부터 초를 계속 세는 것이다.**
유닉스 시간은 1970년 1월 1일 자정(이 날짜는 유닉스가 탄생한 시점이다.) 이후 시간의 초를 사용한다.
이 값을 에폭(epoch)이라 부르며, 에폭은 시스템 간의 날짜와 시간을 교환하는 아주 간단한 방식이다.
에폭 값은 JavaScript와 같은 다른 시스템에서 날짜와 시간을 교환하기 위한 유용한 공통분모다.

## time() 함수
time 모듈의 time() 함수는 현재 시간을 에폭 값으로 반환한다.
```python
import time
now = time.time()
now
...
```
1970년 1월 1일부터 현재까지 10억초가 넘는다.

## ctime() 함수
time 모듈의 ctime() 함수를 사용하여 에폭 값을 문자열로 변환할 수 있다.
```python
import time
now = time.time()
time.ctime(now)
'Fri Apr 5 19:55:32 2019'
```

## struct_time 객체
`time.struct_time` 객체는 `time` 모듈에서 시간을 표현하기 위한 표준적인 데이터 구조이다.
`gmtime()`, `localtime()`, `strptime()`과 같은 여러 `time` 모듈 함수들이 이 `struct_time` 객체를 반환한다.
- **튜플(Tuple)처럼 동작**: `struct_time` 객체는 인덱스 번호를 통해 각 시간 요소에 접근할 수 있다.
- **객체(Object)처럼 동작**: `struct_time` 객체는 점과 속성 이름을 통해 각 시간 요소에 더 명확하게 접근할 수 있다

| 인덱스 | 속성 이름    | 값의 범위           | 설명        |
| --- | -------- | --------------- | --------- |
| 0   | tm_year  | 0000~9999       | 년         |
| 1   | tm_mon   | 1~12            | 월         |
| 2   | tm_mday  | 1~31            | 일         |
| 3   | tm_hour  | 0~23            | 시         |
| 4   | tm_min   | 0~59            | 분         |
| 5   | tm_sec   | 0~59            | 초         |
| 6   | tm_wday  | 0=월요일 ~ 6=일요일   | 요일        |
| 7   | tm_yday  | 1~365           | 년일자       |
| 8   | tm_isdst | 0=아니오 1=예 -1=모름 | 일광 시간 절약제 |

```python
import time
now = time.localtime()
now
time.struct_time(...)

now[0] # 네임드 튜플 처럼 인덱스 사용 가능
2019

now.time_year # 속성 가져오는 것처럼 값 접근 가능
2019
```
## localtime() 메서드
time 모듈의 localtime() 메서드는 시간을 시스템의 표준 시간대로 반환하며, time 모듈의 struct_time 객체로서 반환한다.

> 참고로, 표준 시간대 대신 UTC를 사용하는 것을 추천한다.
> 일광 절약 시간은 사용하지 않는 것이 좋다.

localtime() 메서드는 인자로 에폭 값을 받을 수 있으며, 인자가 주어진 경우에는 주어진 에폭 값을 struct_time 객체로서 반환하게 된다.
```python
import time
now = time.time()
time.localtime(now)
time.struct_time(tm_year...)
```

만약 localtime() 메서드의 인자가 제공되지 않거나 `None`이라면 `time.localtime(time.time())`과 동일하게 동작한다.
```python
import time
time.localtime()
time.struct_time(tm_year...)
```

## gmtime() 메서드
time 모듈의 gmtime() 메서드는 시간을 UTC로 반환하며, time 모듈의 struct_time 객체로서 반환한다.

gmtime() 메서드는 인자로 에폭 값을 받을 수 있으며, 인자가 주어진 경우에는 주어진 에폭 값을 struct_time 객체로서 반환하게 된다.
```python
import time
now = time.time()
time.gmtime(now)
time.struct_time(tm_year...)
```

만약 gmtime() 메서드의 인자가 제공되지 않거나 `None`이라면 `time.localtime(time.time())`과 동일하게 동작한다.
```python
import time
time.gmtime()
time.struct_time(tm_year...)
```

## mktime() 메서드
time 모듈의 mktime() 메서드는 struct_time 객체를 에폭 초로 변환한다.
```python
import time
now = time.time()
tm = time.localtime(now)
time.mktime(tm)
...
```
이 값은 now() 의 에폭 값과 정확하게 일치하지 않는다.
struct_time 객체는 시간을 초까지만 유지하기 때문이다.

## strftime() 메서드, 함수
strftime() 는 날짜와 시간을 특정 포맷 문자열의 형태로 출력할 수 있다.
strftime() 메서드는 datetime, date, time 객체에서 제공된다.
strftime() 함수는 time 모듈에서 제공된다.
strftime() 메서드, 함수는 문자열의 출력 포맷을 지정할 수 있다.
먼저 포맷 문자열 `fmt` 을 정의하고, 이를 사용하자.
숫자는 자릿수에 맞춰 왼쪽에 0이 채워진다.

| 문자열 포맷 | 날짜/시간 단위 | 범위           |
| ------ | -------- | ------------ |
| %Y     | 년        | 1900-...     |
| %m     | 월        | 01-12        |
| %B     | 월 이름     | January, ... |
| %b     | 월 축약 이름  | Jan, ...     |
| %d     | 울의 일자    | 01-31        |
| %A     | 요일 이름    | Sunday, ...  |
| %a     | 요일 축약 이름 | Sun, ...     |
| %H     | 24시간     | 00-23        |
| %I     | 12시간     | 01-12        |
| %p     | 오전/오후    | AM, PM       |
| %M     | 분        | 00-59        |
| %S     | 초        | 00-59        |
이런 특정 포맷을 '형식 지시어'라고 부른다.
참고로, 이 형식 지시어들은 Python에서만 사용하는 것이 아니라, C 언어 라이브러리에서 유래한 **ISO C89 표준**에 기반하고 있다.

time 모듈의 strftime() 함수는 struct_time 객체를 포맷 문자열로서 출력한다.
```python
import time
fmt = "It's %A, %B %d, %Y, local time %I:%M:%S%p"
t = time.localtime()
time.strftime(fmt, t)
"It's Wednesday, March 13, 2019, local time 03:23:46PM"
```

만약 date 객체에 사용하면, 날짜 부분만 변환되게 된다.
그리고 시간은 기본값으로 지정된다.
```python
from datetime import date
some_day = date(2019, 7, 4)
fmt fmt = "It's %A, %B %d, %Y, local time %I:%M:%S%p"
some_day.strftime(fmt)
"It's Thursday, July 04, 2019, local time 12:00:00AM"
```

time 객체는 시간 부분만 변환한다.
```python
from datetime import time
some_time = time(10, 35)
fmt fmt = "It's %A, %B %d, %Y, local time %I:%M:%S%p"
some_time.strftime(fmt)
"It's Monday, January 01, 1900, local time 10:35:00AM"
```

## strptime()
time 모듈의 strptime() 함수를 사용하여 특정 포맷 형태의 문자열을 날짜/시간을 담은 time.struct_time 객체로 반환한다.
특정 포맷 형태의 문자열은 strftime() 에서 사용하는 포맷 문자열이다.
특정 포맷 형태는 사용자가 임의로 지정할 수 있다.
인자로 전달하는 문자열은 특정 포맷 형태에 맞게 전달되어야 한다.
만약 문자열이 특정 포맷 형태가 아닌데 strptime() 함수의 인자로 전달시, 예외가 발생한다.
문자열 포맷은 맞는데, 값 범위가 벗어나면 예외가 발생한다.
```python
import time
fmt = "%Y-%m-%d"
time.strptime("2019-01-19", fmt)
time.struct_time(tm_year...)
```

만약 문자열이 특정 포맷 형태가 아닌데 strptime() 함수의 인자로 전달시, 예외가 발생한다.
```python
import time
fmt = "%Y-%m-%d"
time.strptime("2019 01 19", fmt)
Traceback...
```

위의 경우 특정 포맷 형태를 만약 수정하면 해결된다.
```python
import time
fmt = "%Y %m %d"
time.strptime("2019 01 19", fmt)
time.struct_time(tm_year...)
```

문자열 포맷은 맞는데, 값 범위가 벗어나면 예외가 발생한다.
```python
import time
fmt = "%Y-%m-%d"
time.strptime("2019-13-19", fmt)
Traceback...
```

## 형식 지시어 로케일
| 문자열 포맷 | 날짜/시간 단위 | 범위           |
| ------ | -------- | ------------ |
| %Y     | 년        | 1900-...     |
| %m     | 월        | 01-12        |
| %B     | 월 이름     | January, ... |
| %b     | 월 축약 이름  | Jan, ...     |
| %d     | 울의 일자    | 01-31        |
| %A     | 요일 이름    | Sunday, ...  |
| %a     | 요일 축약 이름 | Sun, ...     |
| %H     | 24시간     | 00-23        |
| %I     | 12시간     | 01-12        |
| %p     | 오전/오후    | AM, PM       |
| %M     | 분        | 00-59        |
| %S     | 초        | 00-59        |
형식 지시어가 만들어내는 결과 문자열은 운영체제의 국제화 설정인 로케일(locale)에 따라 다르다.
다른 월, 일의 이름을 출력하려면 setlocale() 을 사용하여 로케일을 바꿔야 한다.
setlocale() 의 첫 번째 인수는 날짜와 시간을 위한 locale.LC_TIME 이고, 두 번째 인수는 언어와 국가 약어가 결합된 문자열이다.
```python
import locale
from datetime import date
halloween = date(2019, 10, 31)
for lang_country in ['en_us, 'fr_fr', 'de_de', 'es_es', 'is_is']:
	locale.setlocale(locale.LC_TIME, lang_country)
	
	halloween.strftime('%A, %B %d')
```

`lang_country` 값은
```python
import locale
locale.locale_alias.keys()
```
에서 찾을 수 있다.

## 시간 모듈 변환
![시간 상호 변환 요약도](/assets/introduce-python-calc.png)
이는 시간 상호 변환 요약도이다.

## 대체 모듈
- arrow
- dateutil
- iso8601
- fleming
- maya
- dateinfer