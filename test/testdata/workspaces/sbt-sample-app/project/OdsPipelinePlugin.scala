import sbt.Keys.streams
import sbt.io.IO
import sbt.plugins.JUnitXmlReportPlugin.autoImport.testReportsDirectory
import sbt.{AutoPlugin, Setting, Test, file, taskKey}
import scoverage.ScoverageKeys.coverageDataDir
import sbt._

object OdsPipelinePlugin extends AutoPlugin {
  override def trigger = allRequirements

  object autoImport {
    val copyOdsTestReports = taskKey[Unit](
      "copy test reports to the expected ods test report directory (UNIT_TEST_RESULT_DIR)"
    )
    val copyOdsTestCoverageReport = taskKey[Unit](
      "copy test coverage report to the expected location (CODE_COVERAGE_TARGET_FILE)"
    )
    val copyOdsReports =
      taskKey[Unit]("copy all ods reports to the appropriate directories (defined by the envs)")
  }

  import autoImport._

  override lazy val projectSettings: Seq[Setting[_]] = Seq(
    copyOdsTestReports := {
      val log = streams.value.log

      sys.env.get("UNIT_TEST_RESULT_DIR") match {
        case Some(targetPath) => {
          val testReportDir = (Test / testReportsDirectory).value
          log.info(s"copying ${testReportDir.listFiles().length} test report(s) to $targetPath")
          IO.copyDirectory(testReportDir, file(targetPath))
        }
        case None => log.info("no env (UNIT_TEST_RESULT_DIR) set, doing nothing ...")
      }
    },
    copyOdsTestCoverageReport := {
      val log = streams.value.log

      sys.env.get("CODE_COVERAGE_TARGET_FILE") match {
        case Some(targetFilePath) => {
          val scoverageReportFile = coverageDataDir.value / "scoverage-report" / "scoverage.xml"
          log.info(s"copying $scoverageReportFile to $targetFilePath")
          IO.copyFile(scoverageReportFile, file(targetFilePath))
        }
        case None => log.info("no env (CODE_COVERAGE_TARGET_FILE) set, doing nothing ...")
      }
    },
    copyOdsReports := {
      copyOdsTestReports.value
      copyOdsTestCoverageReport.value
    }
  )
}
