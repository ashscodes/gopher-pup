package main
import(."fmt";."math";."strconv")
func d(a,b uint64)uint64{if a>b{return a};return b}
func e(a,b uint64)uint64{if a<b{return a};return b}
func h(x uint64)(n int){for;x>0;x/=10{n++};return}
func g(x uint64)(t uint64){for;x>0;x/=10{t+=1<<(x%10*6)};return}
var t[20]uint64
func init(){t[0]=1;for i:=1;i<20;i++{t[i]=t[i-1]*10}}
func m(x uint64)(f[]uint64){nd:=h(x);if nd&1==1{return};nd/=2;p:=d(t[nd-1],(x+t[nd]-2)/(t[nd]-1));q:=e(x/p,uint64(Sqrt(float64(x))));t:=g(x)
for a:=p;a<=q;a++{b:=x/a;if a*b==x&&(a%10>0||b%10>0)&&t==g(a)+g(b){f=append(f,a)}};return}
func main(){for x,n:=uint64(1),0;x<1e6;x++{if f:=m(x);len(f)>0{n++;Println(FormatUint(x,10))}}}