
// The simplest possible sbt build file is just one line:

scalaVersion := "2.13.1"
// That is, to create a valid sbt build, all you've got to do is define the
// version of Scala you'd like your project to use.

// ============================================================================

// Lines like the above defining `scalaVersion` are called "settings". Settings
// are key/value pairs. In the case of `scalaVersion`, the key is "scalaVersion"
// and the value is "2.13.3"

// It's possible to define many kinds of settings, such as:

name := "martian"
organization := "ch.epfl.scala"
version := "1.0"


libraryDependencies += "org.scala-lang.modules" %% "scala-parser-combinators" % "1.1.2"
libraryDependencies ++= Seq(
  "io.getquill" %% "quill-async-postgres" % "3.5.2", //A
  "org.testcontainers" % "postgresql" % "1.13.0", //B
  "org.postgresql" % "postgresql" % "42.2.11", //C
  "ch.qos.logback"  %  "logback-classic" % "1.2.3" //D
)

