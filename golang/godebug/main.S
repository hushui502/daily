TEXT main.main(SB) D:/project/go/src/awesomeProject2/godebug/main.go
  main.go:4		0x45abe0		65488b0c2528000000	MOVQ GS:0x28, CX			
  main.go:4		0x45abe9		488b8900000000		MOVQ 0(CX), CX				
  main.go:4		0x45abf0		483b6110		CMPQ 0x10(CX), SP			
  main.go:4		0x45abf4		763e			JBE 0x45ac34				
  main.go:4		0x45abf6		4883ec18		SUBQ $0x18, SP				
  main.go:4		0x45abfa		48896c2410		MOVQ BP, 0x10(SP)			
  main.go:4		0x45abff		488d6c2410		LEAQ 0x10(SP), BP			
  main.go:5		0x45ac04		48c744240801000000	MOVQ $0x1, 0x8(SP)			
  main.go:6		0x45ac0d		e8ee32fdff		CALL runtime.printlock(SB)		
  main.go:6		0x45ac12		488b442408		MOVQ 0x8(SP), AX			
  main.go:6		0x45ac17		48890424		MOVQ AX, 0(SP)				
  main.go:6		0x45ac1b		e8b03afdff		CALL runtime.printint(SB)		
  main.go:6		0x45ac20		e87b35fdff		CALL runtime.printnl(SB)		
  main.go:6		0x45ac25		e86633fdff		CALL runtime.printunlock(SB)		
  main.go:7		0x45ac2a		488b6c2410		MOVQ 0x10(SP), BP			
  main.go:7		0x45ac2f		4883c418		ADDQ $0x18, SP				
  main.go:7		0x45ac33		c3			RET					
  main.go:4		0x45ac34		e87773ffff		CALL runtime.morestack_noctxt(SB)	
  main.go:4		0x45ac39		eba5			JMP main.main(SB)			
