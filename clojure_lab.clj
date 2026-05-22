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
