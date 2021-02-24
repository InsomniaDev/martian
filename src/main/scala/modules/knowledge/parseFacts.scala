package modules.knowledge

import modules.databaseQueries._
import io.getquill._

// TODO: Could we possibly export all of this knowledge data into markdown files and then display them through something like Hugo?
// https://gohugo.io/documentation/

class FactParser(ctx: PostgresJdbcContext[SnakeCase.type])
    extends FactData(ctx) {

  private def getNumberOfMatches(value: List[String]): Int = {
    println("hello")

    1
  }

/** getMatchesForAllWords
  * 
  * Get all of the matches for the words that are provided and then order them by count and prevalence
  * 
  * Return the top three results and increment importance
  *
  * @param value
  */
  private def getMatchesForAllWords(value: String) = {
    // TODO: Need to create a unit test for this functionality
    val words = value.split(" ").toList
    val resp = getFactsByUsedWords(words)

    // Need to group by the fact_id
    // **IMPORTANT** This counts on a unique relationship between fact_id and word_id
    // TODO: Add unique relationship between fact_id and word_id in the database
    val distinctFacts = resp.groupBy(_._1.fact_id)

    // Sort by the highest amount of matches
    val highestFoundWords = distinctFacts.toSeq.sortWith(_._1 > _._1)
  }

  /** checkForFact
    *
    * This method retrieves all of the fact data for the provided fact names along with their related facts
    *
    * @param value is the string passed in to look at
    * @return a list of fact data
    */
  private def checkForFact(value: List[String]): List[String] = {

    // Get all of the facts that match the provided list of values
    val resp = checkFactNames(value)

    // TODO: Possibly search for the ones that have the highest count of words...

    // Get all of the fact data from the related facts
    val relatedFacts = getRelatedFactsByIds(
      resp.flatMap(a => a.related_fact_ids.split(";").map(_.toInt))
    ).flatMap(_.fact_data)

    // Take all of the returned facts and put into a single array of facts to return
    resp.flatMap(_.fact_data) ++ relatedFacts
  }

  /** decipherKnowledgeString
    *
    * This method will remove all of the commonly found words from the value passed in and then return the relevant fact set
    *
    * @param value
    * @return the relevant fact set
    */
  def decipherKnowledgeString(value: String): List[String] = {

    // Get all of the commonWords from the database
    val conf = new ConfigData(ctx)
      .getConfigByKey(Some("commonWords"))
      .flatMap(_.value)
      .flatMap(_.split(","))
    println(conf)

    // Get all of the values that aren't in the common words list
    val parsedValues = value.split(" ").view.filter(!conf.contains(_)).toList

    checkForFact(parsedValues)
  }
}
