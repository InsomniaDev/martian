package modules.knowledge

import modules.databaseQueries._
import io.getquill._

// TODO: Could we possibly export all of this knowledge data into markdown files and then display them through something like Hugo?
// https://gohugo.io/documentation/

// TODO: Provide option for the editing user to relate facts
// TODO: Facts should have an importance based on the number of related facts

class FactParser(ctx: PostgresJdbcContext[SnakeCase.type])
    extends FactData(ctx) {

  /** getMatchesForAllWords
    *
    * Get all of the matches for the words that are provided and then order them by count and prevalence
    *
    * Return the top three results and increment importance
    *
    * @param words is a list of words that aren't commonly used and are currently being searched
    */
  private def getMatchesForAllWords(words: List[String]): Seq[FactsToWords] = {

    // TODO: Need to create a unit test for this functionality
    val resp = getFactsByUsedWords(words)

    // Need to group by the fact_id
    // **IMPORTANT** This counts on a unique relationship between fact_id and word_id
    val distinctFacts = resp.groupBy(_._1.fact_id)

    // Sort by the highest amount of matches
    val highestFoundWords =
      distinctFacts.toSeq.sortWith(_._1 > _._1).take(3).flatMap(_._2)

    // Return the matched list sorted by the importance level
    highestFoundWords.map(_._1).sortWith(_.importance > _.importance)
  }

  /** getFactsByFoundIds
    *
    * This will get all of the facts by the provided ids
    * will return a sorted by importance sequence
    *
    * @param factIds
    * @return
    */
  private def getFactsByFoundIds(factIds: Seq[Int]): Seq[Fact] = {

    val foundFacts = getFactsByIds(factIds.toList)

    foundFacts.sortWith(_.importance > _.importance)
  }

  /** decipherKnowledgeString
    *
    * This method will remove all of the commonly found words from the value passed in and then return the relevant fact set
    *
    * @param value
    * @return the relevant fact set
    */
  def decipherKnowledgeString(value: String): Option[List[Fact]] = {

    // Get all of the commonWords from the database
    val conf = new ConfigData(ctx)
      .getConfigByKey(Some("commonWords"))
      .flatMap(_.value)
      .flatMap(_.split(","))
    println(conf)

    // Get all of the values that aren't in the common words list
    val parsedValues = value.split(" ").view.filter(!conf.contains(_)).toList

    // Get all of the top matches by the count of words matched ordered by importance level
    val topImportantMatches = getMatchesForAllWords(parsedValues)

    // If there are relevant matches for the words provided
    topImportantMatches match {
      case tim if (tim.length > 0) => {

        // Increment the importance for the returned matches
        topImportantMatches.foreach(a => incrementFactToWordImportance(a.id))

        // Get the facts for these matches sorted by the fact importance
        Some(getFactsByFoundIds(topImportantMatches.map(_.fact_id)).toList)
      }
      case _ => None
    }
  }
}



// This function might be deprecated
  // private def checkForFact(value: List[String]): List[String] = {

  //   // Get all of the facts that match the provided list of values
  //   val resp = checkFactNames(value)

  //   // TODO: Possibly search for the ones that have the highest count of words...

  //   // Get all of the fact data from the related facts
  //   val relatedFacts = getRelatedFactsByIds(
  //     resp.flatMap(a => a.related_fact_ids.split(";").map(_.toInt))
  //   ).flatMap(_.fact_data)

  //   // Take all of the returned facts and put into a single array of facts to return
  //   resp.flatMap(_.fact_data) ++ relatedFacts
  // }