---
date: 2025-09-15
category: [Arduino]
published: true
fixed: false
---

## arduino-cli 설치
```shell
brew install arduino-cli
```
arduino-cli는 Arduino IDE 보다 더욱 고차원적인 아두이노 관리 프로그램임.

## arduino-cli 최초 설정
```shell
arduino-cli config init
```
arduino-cli를 컴퓨터에 설치한 후 가장 먼저, 딱 한 번만 이 명령어를 실행하면 됨.
이 명령어로 설정 파일의 뼈대를 만든 뒤, ESP8266 보드 URL을 추가하는 등의 설정을 하게 됨.

## WeMos D1 보드 설정 (외부 보드 매니저 주소 추가)
```shell
arduino-cli config set board_manager.additional_urls [http://arduino.esp8266.com/stable/package_esp8266com_index.json](http://arduino.esp8266.com/stable/package_esp8266com_index.json)
```
arduino-cli의 영구적인 설정 파일을 수정하는 명령어임.

arduino-cli의 메인 설정 파일(`arduino-cli.yaml`) 안에 있는 `board_manager.additional_urls` 라는 목록에, 지정된 URL(`http://arduino.esp8266.com/...`)을 추가하고 저장함.

이 명령어는 arduino-cli에게 "앞으로 보드 패키지 목록을 찾을 때, 아두이노 공식 서버뿐만 아니라, 내가 지금 알려주는 이 인터넷 주소도 확인 대상에 포함시켜라" 라고 데이터 소스(Source)를 등록하는 일회성 작업임. 이 명령 자체는 어떠한 다운로드도 수행하지 않음.


## 드라이버 업데이트
```shell
arduino-cli core update-index
```
등록된 모든 데이터 소스로부터 최신 패키지 목록을 다운로드하는 명령어임.

`arduino-cli.yaml` 설정 파일에 등록된 모든 URL(아두이노 공식 URL + 우리가 config set으로 추가한 URL)에 실제로 접속하여, `package_..._index.json` 이라는 이름의 파일들을 다운로드하고 로컬 캐시를 업데이트함.

`package_..._index.json` 파일 안에는 설치 가능한 보드, 라이브러리, 도구들의 이름과 버전 정보가 들어있음.

## WeMos 드라이버 설치
```shell
arduino-cli core install esp8266:esp8266
```
`arduino-cli core update-index` 명령어를 실행 이후에, 최신 목록을 바탕으로 WEMOS D1를 위한 전용 패키지를 설치하기.

## config set ..., core update-index, core install의 차이
- `config set ...` 명령어로 arduino-cli에게 새로운 정보 소스의 주소를 알려줌.

- `core update-index` 명령어는 `config set` 으로 알려준 그 새로운 주소를 포함한 모든 주소에 접속해서, 설치 가능한 패키지 목록을 가져옴.

- `core install ...` 명령어는 `update-index` 가 가져온 그 목록을 보고, 필요한 파일을 설치함.

## Wifi 관련 라이브러리 설치
ESP8266WiFi.h와 같은 WiFi 관련 기능은 별도의 라이브러리를 `lib install`로 설치하는 것이 아님.
이 기능들은 `arduino-cli core install esp8266:esp8266` 명령어로 설치되는 "ESP8266 보드 지원 패키지" 안에 이미 포함되어 있음.

## 아두이노 프로젝트 생성
```shell
arduino-cli sketch new PROJECT_NAME
```

## 아두이노 프로그래밍
`PROJECT_NAME.ino` 이렇게 아두이노 메인 파일 이름은 프로젝트 이름과 동일해야 함이 관례임.

하나의 프로젝트(폴더) 안에는 `setup()`과 `loop()` 함수가 단 한 번만 존재해야 함.
`setup()`과 `loop()`는 프로그램의 시작점임으로, 메인 파일에 위치하는 것이 관례임.

arduino-cli는 프로젝트 내부의 모든 파일들을 독립적으로 인식하는 것이 아니라, 컴파일 직전에 하나의 거대한 파일로 그냥 이어 붙여버림.
컴포넌트 같이 만드러면 `.h` 해더파일을 만들어야함.

아두이노에서 활용하는 언어는 C++의 변형된 프로그래밍 언어를 사용함으로, C++ 프로그래밍과 완전하게 같음.

## 아두이노 관련 VSC 익스텐션 설치
```markdown
Arduino Community Edition
C/C++
C/C++ Extension Pack
```
VSC 익스텐션 3개를 설치해야 함.

## .vscode 디렉토리 생성
`.vscode` 디렉토리를 생성해야 함.
그러나 개발자가 직접 생성 안해도, Arduino 확장 프로그램이 직접 `.vscode` 폴더를 생성하고 `arduino.json` 같은 파일을 거기에 위치시킴.

## VSC에서 아두이노 설정하기
```shell
Arduino: Change Board Type
```
Command Palette에서 `Arduino: Change Board Type` 를 입력한 후, `LOLIN(WeMos) D1 R2 & mini` 선택하기.

이 명령어를 실행하면, `.vscode/arduino.json` 파일이 생성-수정됨.
굳이 arduino.json 파일을 따로 만들지 않아도, .vscode 디렉토리에 자동으로 저장됨.

`arduino.json` 파일은
1. "Arduino" 확장 프로그램의 명령어 생성기 역할
    - Arduino: Change Board Type에서 모델을 선택하면, 그 정보("board": "esp8266:esp8266:d1_mini")가 arduino.json에 저장됨.
    - 이 정보가 VSC의 UI 버튼(Upload(→) 버튼 등)을 누를 때 명령어 인수로서 활용됨.
2. C/C++ 확장 프로그램의 코드 자동완성(IntelliSense)을 위한 정보 제공 역할
을 수행함.

`arduino.json` 파일 내부의 속성은

- "sketch"
    - 어떤 .ino 파일이 메인 파일인지 알려줌.
    - 사용자가 직접 수정해야 할 수 있음. (예: 파일 이름을 바꿨을 때)
- "board
    - 어떤 보드를 타겟으로 컴파일할지 알려줌. (--fqbn 값)
    - Arduino: Change Board Type 메뉴로 설정됨.
- "output"
    - 컴파일 결과물을 어디에 저장할지 알려줌.
    - Arduino: Initialize 실행 시 기본적으로 추가되거나, 사용자가 직접 추가할 수 있음.
- "port"
    - 어느 USB 포트로 업로드할지 알려줌.
    - Arduino: Select Serial Port 메뉴로 설정됨.

이것임.

`arduino.json` 파일은 개발자가 직접 관리할 수 있는 파일임.

## 참고: ESP8266은 같은데 왜 이렇게 Board Type은 선택하는 것이 많은가?
각각 보드가 ESP8266은 이라는 핵심 칩(두뇌)은 동일하지만,
그 두뇌를 둘러싼 '몸체'가 다르기 때문에 VS Code에서 그 많은 모델을 선택해야 하는 것임.

`Arduino: Change Board Type` 에서 특정 모델을 선택하는 것은,
컴파일러에게 바로 이 몸체의 미세한 차이점을 알려주는 매우 중요한 과정임.

## VSC에서 아두이노를 위한 C/C++ IntelliSense 설정하기
```shell
Arduino: Rebuild IntelliSense Configuration
```
`arduino-cli lib install "MQ135"` 이렇게 Arduino 관련 새로운 라이브러리 등을 설치하면,
해당 라이브러리를 `#include` 를 사용해서 불러오면 오류가 발생할 수 있음.
외부 라이브러리를 설치한 후, VS Code가 그 라이브러리의 위치를 인식하도록 강제로 해당 명령어를 활용하여 새로고침하기.

이 명령어를 실행하면, `c_cpp_properties.json` 파일이 생성-수정된다.
`c_cpp_properties.json` 파일은 C/C++ 확장 프로그램의 코드 분석기(IntelliSense)만을 위한 것임.
사용자가 절대로 직접 수정해서는 안 됨.

## 컴파일
```shell
arduino-cli compile --fqbn esp8266:esp8266:d1_mini --output-dir ./build .
```
현재 폴더의 모든 .ino 파일을 합쳐서 컴파일함. 결과물은 build 디렉토리에 저장됨.
`build` 디렉토리는 절대 Git/GitHub에 포함시키면 안됨.

## 업로드
```shell
arduino-cli upload -p /dev/cu.usbserial-A5069RR4 --fqbn esp8266:esp8266:d1_mini .
```
build 디렉토리를 바탕으로 업로드를 수행함.
업로드란 아두이노에게 만든 프로그램을 주입하는 것임.

## 아두이노 보드 검색
```shell
arduino-cli core list
```
해당 명령어로 현재 컴퓨터에 연결된 아두이노 보드 검색하기.
만약 업로드할 때 오류가 난다면, 보드가 정상적으로 연결되지 않았다는 것을 의미함으로, 해당 명령어로 연결된 보드를 검색해보기.

## daemon
```shell
arduino-cli daemon
```
이 명령어를 실행해 두면, arduino-cli를 백그라운드에 상주시켜 컴파일 속도를 비약적으로 향상시킴.

## 출력 확인
```shell
arduino-cli monitor -p /dev/cu.usbserial-A5069RR4 -c baudrate=115200"
```
아두이노 보드가 Serial.print()로 보내는 메시지를 컴퓨터 화면에서 확인할 수 있게하는 명령어.
`-c baudrate` 값은 아두이노 코드의 `Serial.begin()` 에 있는 숫자와 반드시 일치해야 함.
