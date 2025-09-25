; TODO: assert as macro
;(defmacro assert
;  "Evaluates expr and throws an exception if it does not evaluate to
;  logical true."
;  {:added "1.0"}
;  ([x]
;     (when *assert*
;       `(when-not ~x
;          (throw (new AssertionError (str "Assert failed: " (pr-str '~x)))))))
;  ([x message]
;     (when *assert*
;       `(when-not ~x
;          (throw (new AssertionError (str "Assert failed: " ~message "\n" (pr-str '~x))))))))

;; basic eval

(assert true)
(assert (= true true))
(assert (= false false))
(assert (not= false true))
(assert (not false))
(assert (not (not true)))
(assert (= (not false) true))
(assert (= (not true) false))
(assert (not= (not false) false))
(assert (not= (not true) true))
(assert (not nil))
(assert 1234)
(assert (= 1234 1234))
(assert (not= 1234 1235))
(assert "string")
(assert (= "string" "string"))
(assert (not= "string" "strig"))
(assert (not= "string" "strin"))
(assert (not= "string" "stringo"))
(assert [1 2 5])
(assert (= [1 2 5] [1 2 5]))
(assert (not= [1 2 5] [1 2 4]))
(assert (not= [1 2 5] [1 2]))
(assert (not= [1 2 5] [1 2 5 6]))
(assert ())
(assert (= () []))
(assert (not= () '(1)))
; TODO, maps

(assert (= 1234 (inc 1233)))
(assert (= 1234 (+ 1233 1)))
(assert (= (inc 1233) (+ 1233 1)))
(assert (= (inc 1234) (+ 1234 1)))
(assert (not= 1234 (inc 1234)))
(assert (not= 1234 (+ 1234 1)))

; TODO, maps needing eval

;; special forms

; def

(assert (= (def L '(a b c)) 'L))
(assert (= L '(a b c)))
(assert (= (def M 1) 'M))
(assert (= M 1))

; if

(assert (= (if false 3) nil))
(assert (= (if true 3) 3))
(assert (= (if true 3 4) 3))
(assert (= (if (+ 1 2) 3 4) 3))
(assert (= (if nil 3 4) 4))
(assert (= (if () 3 4) 3))
(assert (= (if "" 3 4) 3))
(assert (= (if (if nil nil nil) 3 4) 4))
(assert (= (if (if true nil nil) 3 4) 4))
(assert (= (if (if true true nil) 3 4) 3))
(assert (= (if (if true nil true) 3 4) 4))
(assert (= [1 (if true 3 9) 5] [1 3 5]))

; do

(assert (= (do 1) 1))
(assert (= (do 1 2) 2))
(assert (= (do 1 2 3 true) true))
(assert (= (do 1
               (if true 4 5)
               3
               (if false 9 8))
           8))

; let

(assert (= (let [a 37] a) 37))

; quote

(assert (= '(1 2 3 4 5 6)
           [1 2 3 4 5 6]))
(assert (= (quote 1 2 3 4 5 6)
           1))
(assert (= (quote (1 2 3 4 5 6))
           [1 2 3 4 5 6]))
(assert (= (quote 4) 4))
(assert (= '4 4))
(assert (= (quote A) 'A))
(assert (= 'A 'A))
(assert (= (quote {}) {}))
(assert (= '{} {}))
(assert (= (quote {} {} {}) {}))
(assert (= '({} {} {})
           (quote ({} {} {}))))
(assert (= (quote (+ 1 2))
           '(+ 1 2)))
(assert (= '(+ 1 2)
           '(+ 1 2)))

; fn

(assert (= ((fn [a] [a a]) 3)
           [3 3]))
(assert (= (let [f (fn [a] [a a])]
             (f 4)
             (f 5))
           [5 5]))
(assert (= (let [f (if false 3 (fn [a] [a a]))]
             (f 6)
             (if true (f 33) (f 44)))
           [33 33]))

; arity overrides for fn
; TODO

; loop
; TODO

; recur
; TODO

; var
; TODO

;; destructuring

; vec
(assert (= (let [[x y] [9 10]] x) 9))
(assert (= (let [[x y] [9 10]] y) 10))
(assert (= (let [[x y] '(9 10)] y) 10))
(assert (= (let [[x y] [9]] x) 9))
(assert (= (let [[x y] [9]] y) nil))

; nested let
(assert (= (let [fst (fn [[x y]] x)]
             (let [[x y] ['(1 2 3) '(9 10 11 12)]]
               (fst y)))
           9))
(assert (= (let [scnd (fn [[x y]] y)]
             (let [[x y] ['(1 2 3) '(9 10 11 12)]]
               (scnd y)))
           10))

; nested vector & / "rest"
(assert (= (let [fst (fn [[x y]] x)]
             (let [[x & y] '(9 10 11 12)]
               (fst y)))
           10))
(assert (= (let [scnd (fn [[x y]] y)]
             (let [[x & y] '(9 10 11 12)]
               (scnd y)))
           11))
(assert (= (let [scnd (fn [[x y]] y)]
             (let [[x & y] [9 10 11 12]]
               (scnd y)))
           11))

; multi pair let, with pair
(assert (= (let [fst (fn [[x y]] x)
                 [x y] ['(1 2 3) '(9 10 11 12)]]
             (fst y))
           9))
; multi pair let, with 'rest'
(assert (= (let [scnd (fn [[x y]] y)
                 [x & y] [9 10 11 12]]
             (scnd y))
           11))

; vector binding :as
(assert (= (let [scnd (fn [[x y]] y)]
             (let [[x :as y] '(9 10 11 12)]
               (assert (= x 9))
               (scnd y)))
           10))

; nested bindings, & and :as
(assert (= (let [[[x1 y1][x2 y2]] [[1 2] [3 4]]]
             [x1 y1 x2 y2])
           [1 2 3 4]))
(assert (= (let [[a b & c :as v] [5 6 7 8 9 10]]
             [a b c v])
        [5 6 [7 8 9 10] [5 6 7 8 9 10]]))

; binding string, nested, & and :as, aliasing str
(assert (= (let [[a b & c :as str] "asdjhhfdas"]
             [a b c str])
        [\a \s [\d \j \h \h \f \d \a \s] "asdjhhfdas"]))

; map bindings
;(assert (= (let [{a :a, b :b, c :c, :as m :or {a 2 b 3}}  {:a 5 :c 6}]
;             [a b c m])
;        [5 3 6 {:c 6, :a 5}]))

; nested, map, vec, & and :as
;(assert (= (let [m {:j 15 :k 16 :ivec [22 23 24 25]}
;                 {j :j, k :k, i :i, [r s & t :as v] :ivec, :or {i 12 j 13}} m]
;             [i j k r s t v])
;           [12 15 16 22 23 (24 25) [22 23 24 25]]))

;:keys, :strs and :syms
;prefixed keys
;(let [m {:x/a 1, :y/b 2}
;      {:keys [x/a y/b]} m]
;  (+ a b))
;-> 3
;bound keys
;(let [m {::x 42}
;      {:keys [::x]} m]
;  x)
;-> 42

; aliasing
(assert (= (let [x 7]
             (let [x (inc x)]
               (assert (= x 8)))
             x)
           7))

; vector :as
; TODO

; map
; TODO

; some math

(assert (= (+ 4) 4))
(assert (= (+ 3 4) 7))
(assert (= (+ 2 3 4) 9))
(assert (= (* 4) 4))
(assert (= (* 3 4) 12))
(assert (= (* 2 3 4) 24))

; misc clojure

;(package main)
;(import fmt)

(def kwvalue :keyword)
(def mapvalue {1 "abc" :key :val false true :nil nil})
(def nilvalue nil)

;(defn fact [n] (if (<= n 1) 1 (* n (fact (- n 1)))))

;(defn main []
;  (fmt/Println "hello world"))

;; map, name

; expect -> (a b c)
;(map name L)
; expect -> ("a" "b" "c")

;; list building
;(cons 42 (cons 69 (cons 613 nil)))
; expect -> (42 69 613)
;(list 42 69 613)
; expect -> (42 69 613)

;; list walking, first, list, second, rest

;(first (list 1 2))
; expect -> 1
;(second (list 1 2))
; expect -> 2
;(second '(1 2 3))
; expect -> 2
;(first (rest '(1 2 3)))
; expect -> 2
;(first (first '((1 2) (3 4))))
; expect -> 1

;; string
;(show "I'm string")

;; hash from two arrays
;(let (Keys '(one two three)  Values (1 2 3))
;   (mapc println
;      (mapcar cons Keys Values) ) )

; expect -> (one . 1)
; expect -> (two . 2)
; expect -> (three . 3)

;;; read macros
;'(a b c)
;; expect -> (a b c)
;'(quote . a)
;; expect -> 'a
;;(cons 'quote 'a)   # (quote . a)
;; expect -> 'a
;(list 'quote 'a)   # (quote a)
;; expect -> '(a)

;; read macros; backtick

;'(a `(+ 1 2 3) z)
;; expect -> (a 6 z)

;;; read macros; tilde

;'(a b c ~(list 'd 'e 'f) g h i)
;; expect -> (a b c d e f g h i)

;;; funarg test
;(define (adder n)
;   (lambda (x) (+ x n)))

;(define add5 (adder 5))
;(add5 4)
;; expect -> 9

;(let* ((n 5)
;       (add5 (adder n)))
;   (let ((n 10))
;      (add5 4)))
; expect -> 14

;(let (N 7  L (7 7 7)) (inc 'N) (inc (cdr L)) (cons N L))
; expect -> (8 7 8 7)

;(mapcar if '(true nil true nil) '(1 2 3 4) '(5 6 7 8))
; expect -> (1 6 3 8)

;{ I'm
; a block comment!
; }
;}#
