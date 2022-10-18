addSbtPlugin("com.github.sbt" % "sbt-native-packager" % "1.9.4")
addSbtPlugin("org.scalameta"  % "sbt-scalafmt"        % "2.4.6")
addSbtPlugin("org.scoverage"  % "sbt-scoverage"       % "2.0.5")

// this fixes the problem with different versions of scala-xml in twirl and the scoverage sbt plugin :F
libraryDependencySchemes += "org.scala-lang.modules" %% "scala-xml" % VersionScheme.Always
