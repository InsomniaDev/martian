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

/**
  * checkForFact
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

/**
  * decipherKnowledgeString
  * 
  * This method will remove all of the commonly found words from the value passed in and then return the relevant fact set
  *
  * @param value
  * @return the relevant fact set
  */
  def decipherKnowledgeString(value: String): List[String] = {

    // Get all of the commonWords from the database
    val conf = new ConfigData(ctx).getConfigByKey(Some("commonWords")).flatMap(_.value).flatMap(_.split(","))
    println(conf)

    // Get all of the values that aren't in the common words list
    val parsedValues = value.split(" ").view.filter(!conf.contains(_)).toList

    checkForFact(parsedValues)
  }
}
