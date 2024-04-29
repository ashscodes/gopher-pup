package main
import (. "fmt";. "os";. "strconv";. "strings")
func main(){b:=func(n int)(s string){v,r:=[]int{1e3,900,500,400,100,90,50,40,10,9,5,4,1},Fields("M CM D CD C XC L XL X IX V IV I");for i:=0;n>0;i++{for v[i]<=n{s+=r[i];n-=v[i]}};return};for _,a:=range Args[1:]{n,_:=Atoi(a);Println(b(n))}}