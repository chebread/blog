---
desc: 해커와 화가 (폴 그레이엄)
date: 2025-09-13
category: [독서]
published: true
fixed: false
---

> 또 전체적인 프로그램을 미리 신중하게 적어서 생각하는 방향이 옳은지 여부를 확인하기 전에 조각난 코드부터 대책 없이 늘어놓은 다음 그것의 모양을 조금씩 잡아 나가는  방법으로 프로그래밍을 했다.  
> 소설과, 화가, 그리고 건축가의 직업이 그런 것처럼 프로그램이란 전체 모습을 미리 알 수 있는 것이 아니라 작성해 나가면서 이해하게 되는 존재다.  
> 프로그래밍 언어는 당신이 이미 머릿속으로 생각한 프로그램을 표현하는 도구가 아니라, 아직 존재하지 않는 프로그램을 생각해 내기 위한 도구다.  
> 해커에게 필요한 언어는 마음껏 내갈기고, 더럽히고, 사방에 떡칠할 수 있는 언어다.  
> — 해커와 화가

그렇다. 코드는 더러워야 한다. 정결한 코드는 아무 의미가 없다. 코드는 더러움 속에서, 혼란 속에서 존재되어야만 한다. 코드에 아이디어를 떡칠해야 한다. 왜 코드 외적인 것에서 아이디어를 피력하고, 생산하는가? 항상 코드 자체에 아이디어를 적어 놓아야 한다. 프로그래머란, 오로지 코드에 어떤 것을 피력해야 한다. 코드 외적인 것에서 놀지 마라. 코드에서 놀아라. 코드는 깔끔해야 한다는 강박을 버려라. 코드는 더러워야 한다.

```c
// https://en.wikipedia.org/wiki/Fast_inverse_square_root
float Q_rsqrt( float number )
{
	long i;
	float x2, y;
	const float threehalfs = 1.5F;

	x2 = number * 0.5F;
	y  = number;
	i  = * ( long * ) &y;                       // evil floating point bit level hacking
	i  = 0x5f3759df - ( i >> 1 );               // what the fuck?
	y  = * ( float * ) &i;
	y  = y * ( threehalfs - ( x2 * y * y ) );   // 1st iteration
//	y  = y * ( threehalfs - ( x2 * y * y ) );   // 2nd iteration, this can be removed

	return y;
}
```

오픈소스로 배포한다고 해도 그냥 더러운 코드를 공개하라. 왜 당신의 아이디어가 담긴 더러운 코드는 공개하기 싫은가?
프로그래밍의 거장도 더럽게 코드를 짜고, 더러운 코드를 공개하는데 당신은 왜 깨끗하려 하는가?

> Chaos to Cosmos.