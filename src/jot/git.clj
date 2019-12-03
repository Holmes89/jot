(ns jot.git
  (:require [clojure.java.shell :refer [sh]]))

(def jot-base
  (str (System/getProperty "user.home") "/.jot"))

(defn pull
  []
  (sh "git" "-C" jot-base "pull" "origin" "master"))

(defn push
  []
  (sh "git" "-C" jot-base "push" "origin" "master"))

(defn commit
  [msg]
  (sh "git" "-C" jot-base "commit" "-m" msg))

(defn add
  [filename]
  (sh "git" "-C" jot-base "add" filename))
