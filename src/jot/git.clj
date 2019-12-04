(ns jot.git
  (:require [clojure.java.shell :refer [sh]]))

(def jot-base
  (str (System/getProperty "user.home") "/.jot"))

(defn success?
  [result-map]
  (= 0 (:exit result-map)))

(defn pull
  []
  (success? (sh "git" "-C" jot-base "pull" "origin" "master")))

(defn push
  []
  (success? (sh "git" "-C" jot-base "push" "origin" "master")))

(defn commit
  [msg]
  (success? (sh "git" "-C" jot-base "commit" "-m" msg)))

(defn add
  [filename]
  (success? (sh "git" "-C" jot-base "add" filename)))
