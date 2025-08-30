![giveMeSub Logo](assets/logo.png)

# giveMeSub
> Discovering the surface at the speed of Go.

`giveMeSub` - bu Go tilida yozilgan yuqori tezlikdagi, parallel ishlaydigan subdomenlarni aniqlash vositasi. U berilgan domen uchun subdomenlarni topadi, ularning IP manzillarini aniqlaydi va web-serverlar ishlayotganini tekshiradi. Loyiha muallifi: **GradientSec @ethica**.

---

## 🚀 Asosiy Xususiyatlari

- **Yuqori Tezlik:** Go'ning `goroutine`lari yordamida bir vaqtning o'zida yuzlab so'rovlarni jo'natish.
- **DNS & Web-Server Tekshiruvi:** Nafaqat subdomen mavjudligini, balki uning `HTTP` (80) va `HTTPS` (443) portlari ochiqligini ham tekshiradi.
- **Faylga Saqlash:** Barcha topilgan natijalarni toza formatda `.txt` fayliga avtomatik saqlaydi.
- **Moslashuvchanlik:** Parametrlar (`-d`, `-w`, `-o`, `-t`) orqali oson boshqaruv.

## ⚙️ O'rnatish

1.  GitHub sahifasining "Releases" bo'limidan o'zingizning operatsion tizimingiz uchun tayyor dasturni yuklab oling.
2.  Yoki manba kodidan o'zingiz kompilyatsiya qiling:
    ```bash
    git clone [https://github.com/SizningProfilingiz/giveMeSub.git](https://github.com/SizningProfilingiz/giveMeSub.git)
    cd giveMeSub
    go build -o giveMeSub
    ```

## 📋 Ishlatish

Dasturni terminal orqali quyidagi parametrlar bilan ishga tushiring:

```bash
./giveMeSub -d <maqsad_domen> -w <wordlist_fayli> -o <natija_fayli> -t <potoklar_soni>
```

**Namuna:**
```bash
./giveMeSub -d example.com -w subdomains.txt -o found.txt -t 200
```

## 📊 Chiqish Namuna (Output Example)

```
[+] Topildi: [www.example.com](https://www.example.com)              -> IP: [93.184.216.34]    -> Portlar: [HTTP:80, HTTPS:443]
[+] Topildi: mail.example.com              -> IP: [192.0.2.1]        -> Portlar: [N/A]
...
```

## 📝 Litsenziya

Ushbu loyiha MIT Litsenziyasi ostida tarqatiladi. To'liq ma'lumot uchun `LICENSE` faylini ko'ring.
