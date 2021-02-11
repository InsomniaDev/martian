package modules.knowledge

import modules.databaseQueries._
import io.getquill._

// TODO: Could we possibly export all of this knowledge data into markdown files and then display them through something like Hugo?
// https://gohugo.io/documentation/

class FactParser(ctx: PostgresJdbcContext[SnakeCase.type])
    extends FactData(ctx) {
  // TODO: Need to check for fact existence here
  // TODO: Need to use config values in here and have a class with those values

  // TODO: We should look into a type of sentence parser here

  // checkForFact
  //
  // @param value is the offered string
  private def checkForFact(value: String) = {

    // TODO: Make this method have a List[String] param that will have a list of facts that don't include the most common words

    // Get all of the facts that match the provided name
    val resp = checkFactNames(value.split(" ").toList)

    // Get all of the fact data from the related facts
    val relatedFacts = getRelatedFactsByIds(
      resp.flatMap(a => a.related_fact_ids.split(";").map(_.toInt))
    ).flatMap(_.fact_data)

    // Take all of the returned facts and put into a single array of facts to return
    val allFactData = resp.flatMap(_.fact_data) ++ relatedFacts
  }
}
