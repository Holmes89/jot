(ns jot.entry
  (:require [clojure.java.shell :refer [sh]]
            [clojure.java.io :as io]
            [clojure.string :refer [join]]
            [jot.git :as git]))

(def editor "emacs")

(defn temp-file-name
  []
  (str "/tmp/jot-entry-" (rand-int 100000) ".md"))

(defn open-editor
  [name]
  (sh editor name))

(defn date [] (java.util.Date.))

(def current-entry-date-string
  (.format (java.text.SimpleDateFormat. "yyyy-MM-dd") (date)))

(def file-header
  (str "# " (.format (java.text.SimpleDateFormat. "EEEE, MMMM d, yyyy") (date)) "\n\n"))

(defn entry-name
  ([] (entry-name current-entry-date-string))
  ([date-string] (str git/jot-base "/entries/" date-string ".md")))

(defn current-entry-timestamp-string
  []
  (.format (java.text.SimpleDateFormat. "HH:mm:ss zz") (date)))

(defn entry-header
  []
  (str "## " (current-entry-timestamp-string)))

(defn file-exists?
  [filename]
  (.exists (io/as-file filename)))

(defn append-entry
  [contents]
  (let [file (entry-name)]
    (if-not (file-exists? file)
      (spit file file-header))
    (spit file (join "\n" [(entry-header) "" contents ""]) :append true)))

(defn create
  []
  (let [name (temp-file-name)]
    (cond
      (not (git/pull))
      {:exit 1 :message "unable to pull from repo"}

      (not (open-editor name))
      {:exit 1 :message "failed to open editor"}

      (not (append-entry (slurp name)))
      {:exit 1 :message "failed to write to entry"}
      )    
    (io/delete-file name)
    (git/add (entry-name))
    (git/commit (str "entry created " (current-entry-timestamp-string)))
    (git/push)
    true)) ;; TODO error handling?


(defn show
  [date]
  (let [file-name (entry-name date)]
    (if (file-exists? file-name)
      (println (slurp file-name))
      (println "entry does not exist")))
  true)
