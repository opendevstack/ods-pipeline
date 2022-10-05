name         := """sbt-sample-app"""
organization := "org.opendevstack.pipeline.sbt"

version := "1.0-SNAPSHOT"

lazy val root = (project in file(".")).enablePlugins(PlayScala)

scalaVersion := "2.13.9"

// this fixes the problem with different versions of scala-xml in twirl and the scoverage sbt plugin :F
libraryDependencySchemes += "org.scala-lang.modules" %% "scala-xml" % VersionScheme.Always

libraryDependencies += guice
libraryDependencies += "org.scalatestplus.play" %% "scalatestplus-play" % "5.0.0" % Test

// Adds additional packages into Twirl
//TwirlKeys.templateImports += "org.opendevstack.pipeline.sbt.controllers._"

// Adds additional packages into conf/routes
// play.sbt.routes.RoutesKeys.routesImport += "org.opendevstack.pipeline.sbt.binders._"
