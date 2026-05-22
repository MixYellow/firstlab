  (:require [clojure.string :as str]))

;; очистка
(defn clean-str [s]
  (let [trimmed (str/trim s)]
    (when-not (empty? trimmed) trimmed)))

;; трансформация
(defn str->number [s]
  (try
    (if (str/includes? s ".")
      (Double/parseDouble s)
      (Long/parseLong s))
    (catch NumberFormatException _
      nil)))

;; валидация 
(defn positive? [x]
  (and (number? x) (> x 0)))

(defn process-element [s]
  (some-> s
          clean-str
          str->number
          (#(when (positive? %) %))))

(defn pipeline [raw-strings]
  (->> raw-strings
       (map process-element)
       (filter some?)
       (into [])))
