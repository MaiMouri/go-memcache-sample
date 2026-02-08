# Go Cache Library Comparison

Goで利用可能なキャッシュライブラリの比較リポジトリ。

## ライブラリ一覧

| # | ライブラリ | パッケージ | 種類 | GitHub Stars |
|---|-----------|-----------|------|-------------|
| 1 | gomemcache | `github.com/bradfitz/gomemcache` | Memcachedクライアント | - |
| 2 | rainycape | `github.com/rainycape/memcache` | Memcachedクライアント | - |
| 3 | go-cache | `github.com/patrickmn/go-cache` | インメモリ | - |
| 4 | ristretto | `github.com/dgraph-io/ristretto` | インメモリ | - |
| 5 | freecache | `github.com/coocood/freecache` | インメモリ | - |

## 特徴比較

### Memcachedクライアント

| 特徴 | gomemcache | rainycape |
|------|-----------|-----------|
| 外部サーバー | 必要 (memcached) | 必要 (memcached) |
| プロトコル | テキスト/バイナリ | テキスト/バイナリ |
| コネクションプール | なし | あり |
| 値の型 | `[]byte` | `[]byte` |
| TTL | サーバー側で管理 (秒単位) | サーバー側で管理 (秒単位) |
| 分散キャッシュ | 対応 (複数サーバー指定可) | 対応 (複数サーバー指定可) |
| メンテナンス状況 | 活発 (Brad Fitzpatrick作) | 非活発 (最終更新2015年) |

### インメモリキャッシュ

| 特徴 | go-cache | ristretto | freecache |
|------|----------|-----------|-----------|
| 外部サーバー | 不要 | 不要 | 不要 |
| 値の型 | `interface{}` (任意の型) | `interface{}` (任意の型) | `[]byte` |
| TTL | あり (`time.Duration`指定) | あり (`SetWithTTL`) | あり (秒単位) |
| エビクションポリシー | TTLベース | TinyLFU (アドミッション + SampledLFU) | LRUベース |
| メモリ上限指定 | なし | あり (MaxCost) | あり (初期化時にサイズ指定) |
| GCへの影響 | 大 (mapベース) | 中 | ゼロ (ポインタ不使用) |
| 並行安全性 | あり (sync.RWMutex) | あり (ロックフリー設計) | あり (セグメントロック) |
| Setの即時反映 | 即時 | 非同期 (バッファ経由) | 即時 |
| メトリクス | なし | あり (Hits/Misses) | あり (HitRate/EntryCount) |
| 適したユースケース | シンプルなキャッシュ、小〜中規模 | 高スループット、大規模、ヒット率最適化 | 大量エントリ、GC回避が必要な環境 |

## 選定ガイド

### Memcachedクライアントが必要な場合

- **gomemcache**: メンテナンスが活発でデファクトスタンダード。特別な理由がなければこちらを選択。
- **rainycape**: コネクションプールが組み込みで必要な場合。ただしメンテナンスが停滞。

### インメモリキャッシュが必要な場合

- **go-cache**: シンプルなAPIで導入が容易。小〜中規模のキャッシュに最適。任意のGoの型をそのまま格納できる。
- **ristretto**: 高スループットが必要な場合。TinyLFUによりキャッシュヒット率が高い。ただしSetが非同期のため即時反映が必要な場合は注意。
- **freecache**: 大量のキャッシュエントリがある環境でGCの影響を最小化したい場合。値は`[]byte`のみのためシリアライズが必要。

## ディレクトリ構成

```
.
├── README.md                          # このファイル
├── benchmark_test.go                  # 全ライブラリのベンチマーク
├── examples/
│   ├── gomemcache/main.go             # gomemcache使用例
│   ├── rainycape/main.go              # rainycape使用例
│   ├── gocache/main.go                # go-cache使用例
│   ├── ristretto/main.go              # ristretto使用例
│   └── freecache/main.go              # freecache使用例
├── go.mod
├── go.sum
└── main.go
```

## ベンチマークの実行

### インメモリキャッシュのみ (memcached不要)

```bash
go test -bench='GoCache|Ristretto|Freecache' -benchmem -count=3
```

### 全ライブラリ (memcachedサーバーが必要)

```bash
# memcachedを起動
# macOS: brew install memcached && memcached -d
# Docker: docker run -d -p 11211:11211 memcached

go test -bench=. -benchmem -count=3
```

### ベンチマーク一覧

| ベンチマーク名 | 内容 |
|---------------|------|
| `BenchmarkGomemcache` | gomemcache Set+Get (memcached必要) |
| `BenchmarkRainycape` | rainycape Set+Get (memcached必要) |
| `BenchmarkGoCache` | go-cache Set+Get |
| `BenchmarkRistretto` | ristretto Set+Get (非同期) |
| `BenchmarkRistrettoWithWait` | ristretto Set+Get (反映待ちあり) |
| `BenchmarkFreecache` | freecache Set+Get |
| `BenchmarkGoCacheParallel` | go-cache 並行Set+Get |
| `BenchmarkRistrettoParallel` | ristretto 並行Set+Get |
| `BenchmarkFreecacheParallel` | freecache 並行Set+Get |
| `BenchmarkGoCacheManyKeys` | go-cache ユニークキーSet+Get |
| `BenchmarkRistrettoManyKeys` | ristretto ユニークキーSet+Get |
| `BenchmarkFreecacheManyKeys` | freecache ユニークキーSet+Get |

## 各ライブラリの詳細

### 1. gomemcache (`bradfitz/gomemcache`)

Memcachedの公式Goクライアント。Brad Fitzpatrick (memcachedの作者) によるもの。

```go
mc := memcache.New("127.0.0.1:11211")
mc.Set(&memcache.Item{Key: "foo", Value: []byte("bar"), Expiration: 300})
item, err := mc.Get("foo")
```

### 2. rainycape (`rainycape/memcache`)

gomemcacheの代替Memcachedクライアント。コネクションプールを内蔵。

```go
mc, err := memcache.New("127.0.0.1:11211")
mc.Set(&memcache.Item{Key: "foo", Value: []byte("bar"), Expiration: 300})
item, err := mc.Get("foo")
```

### 3. go-cache (`patrickmn/go-cache`)

シンプルなインメモリKey-Valueキャッシュ。`interface{}`で任意の型を格納可能。

```go
c := cache.New(5*time.Minute, 10*time.Minute)
c.Set("foo", "bar", cache.DefaultExpiration)
val, found := c.Get("foo")
```

### 4. ristretto (`dgraph-io/ristretto`)

Dgraph社製の高性能インメモリキャッシュ。TinyLFUアドミッションポリシーにより高いヒット率を実現。

```go
cache, _ := ristretto.NewCache(&ristretto.Config{
    NumCounters: 1e7,
    MaxCost:     1 << 30,
    BufferItems: 64,
})
cache.Set("foo", "bar", 1)
time.Sleep(10 * time.Millisecond) // 非同期反映を待つ
val, found := cache.Get("foo")
```

**注意**: Setは非同期でバッファ経由で処理されるため、Set直後のGetではまだ反映されていない場合がある。

### 5. freecache (`coocood/freecache`)

GCオーバーヘッドゼロのインメモリキャッシュ。内部的にポインタを使用しないことでGCスキャンを回避。

```go
cache := freecache.NewCache(100 * 1024 * 1024) // 100MB
cache.Set([]byte("foo"), []byte("bar"), 300)    // TTL: 300秒
val, err := cache.Get([]byte("foo"))
```

**注意**: キーと値は`[]byte`のみ。構造体等を格納するにはシリアライズ (JSON, gob等) が必要。
