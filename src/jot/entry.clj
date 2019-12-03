(ns jot.editor
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

(def entry-name
  (str git/jot-base "/entries/" current-entry-date-string ".md"))

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
  (let [file entry-name]
    (if-not (file-exists? file)
      (spit file file-header))
    (spit file (join "\n" [(entry-header) "" contents ""]) :append true)))

(defn create-entry
  []
  (let [name (temp-file-name)]
    (git/pull)
    (open-editor name)
    (append-entry (slurp name))
    (io/delete-file name)
    (git/add entry-name)
    (git/commit (str "entry created " (current-entry-timestamp-string)))
    (git/push)))
