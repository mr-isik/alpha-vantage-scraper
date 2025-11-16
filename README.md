# Alpha Vantage Scraper

Bu, Alpha Vantage API'sinden veri çekmek ve CSV dosyası olarak kaydetmek için Go ile yazılmış bir CLI (Komut Satırı Arayüzü) uygulamasıdır. Hem hisse senedi piyasası verilerini hem de haberleri çekebilirsiniz.

## Özellikler

- Geçmiş hisse senedi verilerini çekme (TIME_SERIES_DAILY)
- Haberleri ve duyarlılık analizini çekme (NEWS_SENTIMENT)
- Verileri tarih aralığına göre filtreleme
- Çıktıyı CSV dosyası olarak kaydetme
- `.env` dosyası üzerinden API anahtarı yönetimi

## Kurulum (Windows İçin)

1.  **Depoyu klonlayın:**

    ```powershell
    git clone https://github.com/mr-isik/alpha-vantage-scraper.git
    cd alpha-vantage-scraper
    ```

2.  **Bağımlılıkları yükleyin:**

    ```powershell
    go mod tidy
    ```

3.  **API anahtarınızı ayarlayın:**

    Projenin kök dizininde bir `.env` dosyası oluşturun ve Alpha Vantage API anahtarınızı ekleyin:

    ```
    API_KEY=API_ANAHTARINIZ
    ```

    Ücretsiz bir API anahtarını [Alpha Vantage web sitesinden](https://www.alphavantage.co/support/#api-key) alabilirsiniz.

4.  **Uygulamayı derleyin:**
    Uygulamayı çalıştırılabilir bir `.exe` dosyası haline getirin.

    ```powershell
    go build -o av-scraper.exe
    ```

## Detaylı Kullanım Rehberi (Windows İçin)

**Önemli Not:** Bu bir komut satırı uygulamasıdır. `av-scraper.exe` dosyasına çift tıklayarak **çalıştırılamaz**. Komutları, projenin bulunduğu klasörde açtığınız bir terminal (PowerShell veya Komut İstemi) üzerinden çalıştırmanız gerekir.

Uygulama, `stocks` (hisse senetleri) ve `news` (haberler) olmak üzere iki ana komut üzerinden çalışır.

---

### **1. Hisse Senedi Verilerini Çekme (`stocks`)**

Bu komut, belirli bir hisse senedi için günlük zaman serisi verilerini (açılış, en yüksek, en düşük, kapanış, hacim) çeker.

**Komut Yapısı:**

```powershell
./av-scraper.exe stocks --symbol=<SEMBOL> [diğer bayraklar]
```

**Bayraklar (Flags):**

| Bayrak         | Açıklama                                                                | Gerekli mi? | Varsayılan Değer |
| -------------- | ----------------------------------------------------------------------- | ----------- | ---------------- |
| `--symbol`     | Verisi çekilecek hisse senedinin sembolü (örn: `IBM`, `MSFT`, `GOOGL`). | **Evet**    | -                |
| `--start-date` | Veri çekilecek başlangıç tarihi (`YYYY-AA-GG` formatında).              | Hayır       | Tüm veriler      |
| `--end-date`   | Veri çekilecek bitiş tarihi (`YYYY-AA-GG` formatında).                  | Hayır       | Tüm veriler      |
| `--output`     | Sonuçların kaydedileceği CSV dosyasının adı.                            | Hayır       | `output.csv`     |

**Örnek Kullanımlar:**

- **Tüm Geçmiş Verileri Çekme:**
  Microsoft (`MSFT`) için mevcut olan tüm geçmiş verileri çeker ve varsayılan olarak `output.csv` dosyasına kaydeder.

  ```powershell
  ./av-scraper.exe stocks --symbol=MSFT
  ```

- **Belirli Bir Tarih Aralığı İçin Veri Çekme:**
  Apple (`AAPL`) için 1 Haziran 2023 ile 30 Haziran 2023 arasındaki verileri çeker.

  ```powershell
  ./av-scraper.exe stocks --symbol=AAPL --start-date=2023-06-01 --end-date=2023-06-30
  ```

- **Özel Çıktı Dosyası Belirleme:**
  Tesla (`TSLA`) için tüm verileri çeker ve `tesla_verileri.csv` adında bir dosyaya kaydeder.

  ```powershell
  ./av-scraper.exe stocks --symbol=TSLA --output=tesla_verileri.csv
  ```

- **Tek Bir Günün Verisini Çekme:**
  Google (`GOOGL`) için sadece 15 Kasım 2023 tarihinin verisini almak için başlangıç ve bitiş tarihini aynı girin.
  ```powershell
  ./av-scraper.exe stocks --symbol=GOOGL --start-date=2023-11-15 --end-date=2023-11-15 --output=google_15kasim.csv
  ```

---

### **2. Haber ve Duyarlılık Analizi Verilerini Çekme (`news`)**

Bu komut, belirli hisse senetleri, konular veya genel piyasa hakkında haberleri ve bu haberlerin duyarlılık analizini (pozitif, negatif, nötr) çeker.

**Komut Yapısı:**

```powershell
./av-scraper.exe news [bayraklar]
```

**Bayraklar (Flags):**

| Bayrak         | Açıklama                                                                                 | Gerekli mi?                | Varsayılan Değer |
| -------------- | ---------------------------------------------------------------------------------------- | -------------------------- | ---------------- |
| `--tickers`    | Haberleri aranacak hisse senedi sembolleri (virgülle ayrılmış, örn: `IBM,AAPL`).         | Hayır (en az biri gerekli) | -                |
| `--topics`     | Haberleri aranacak konular (virgülle ayrılmış, örn: `technology,finance`, `blockchain`). | Hayır (en az biri gerekli) | -                |
| `--start-date` | Haberler için başlangıç zamanı (`YYYYAAGGTSSDD` formatında, örn: `20230101T0000`).       | Hayır                      | Tüm zamanlar     |
| `--end-date`   | Haberler için bitiş zamanı (`YYYYAAGGTSSDD` formatında, örn: `20230131T2359`).           | Hayır                      | Tüm zamanlar     |
| `--limit`      | Döndürülecek maksimum haber sayısı.                                                      | Hayır                      | `50`             |
| `--output`     | Sonuçların kaydedileceği CSV dosyasının adı.                                             | Hayır                      | `output.csv`     |

**Önemli Not:** `tickers` veya `topics` bayraklarından en az birini sağlamanız gerekmektedir.

**Örnek Kullanımlar:**

- **Belirli Hisseler İçin Haber Çekme:**
  IBM ve Apple (`IBM,AAPL`) ile ilgili en son 50 haberi çeker.

  ```powershell
  ./av-scraper.exe news --tickers=IBM,AAPL
  ```

- **Belirli Bir Konu Hakkında Haber Çekme:**
  `blockchain` konusuyla ilgili en son 50 haberi çeker ve `blockchain_haberleri.csv` dosyasına kaydeder.

  ```powershell
  ./av-scraper.exe news --topics=blockchain --output=blockchain_haberleri.csv
  ```

- **Birden Fazla Konu ve Hisse İçin Kapsamlı Arama:**
  `mergers_and_acquisitions` (birleşme ve devralmalar) konusuyla ilgili olarak `MSFT` ve `GOOGL` hisselerini içeren haberleri arar.

  ```powershell
  ./av-scraper.exe news --tickers=MSFT,GOOGL --topics=mergers_and_acquisitions
  ```

- **Zaman Aralığı ve Limit Belirterek Haber Çekme:**
  `earnings` (kazanç raporları) konusuyla ilgili, 2023'ün ilk çeyreğinde yayınlanmış en son 200 haberi çeker ve `kazanc_raporlari_2023Q1.csv` dosyasına kaydeder.
  ```powershell
  ./av-scraper.exe news --topics=earnings --start-date=20230101T0000 --end-date=20230331T2359 --limit=200 --output=kazanc_raporlari_2023Q1.csv
  ```

## Katkıda Bulunma

Pull request'ler kabul edilir. Büyük değişiklikler için, lütfen önce neyi değiştirmek istediğinizi tartışmak üzere bir issue açın.

## Lisans

[MIT](https://choosealicense.com/licenses/mit/)
