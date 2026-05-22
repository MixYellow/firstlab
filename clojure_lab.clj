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

(def sample-data
  ["  123  "
   " -5   "
   "  45.6 "
   "  abc  "
   "   0   "
   "   789 "
   "   "
   " 12.34 "])

(println "Result:")
(println (pipeline sample-data))

(defn read-lines-until-empty []
  (println "Enter lines (press Enter on empty line to finish):")
  (loop [lines []]
    (let [line (read-line)]
      (if (or (nil? line) (empty? line))
        lines
        (recur (conj lines line))))))

(defn -main []
  (println "=== Data Processing Pipeline ===")
  (let [user-input (read-lines-until-empty)]
    (println "\nYou entered:")
    (println user-input)
    (println "\nPipeline result:")
    (println (pipeline user-input))))

(-main)
