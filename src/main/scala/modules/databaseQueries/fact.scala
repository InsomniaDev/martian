// package modules.databaseQueries

// import io.getquill._
// import java.util.UUID
// import scala.concurrent.{Future, ExecutionContext}
// // val q = quote {
// //   query[Book].filter(p => p.pages.contains(25)).allowFiltering
// // }
// // ctx.run(q)
// // record_uuid uuid,
// // account_uuid uuid,
// // tags set<text>,
// // words set<text>,
// // record text,
// // importance int,

// case class Records(
//     record_uuid: UUID,
//     account_uuid: UUID,
//     tags: Set[String],
//     words: Set[String],
//     record: String,
//     importance: Int
// )

// class FactData(ctx: CassandraSyncContext[SnakeCase.type])(implicit
//     ec: ExecutionContext
// ) {
//   import ctx._

//   /** getFactsByUsedWord
//     *
//     * SELECT
//     *   record_uuid,
//     *   account_uuid,
//     *   tags,
//     *   words,
//     *   record,
//     *   importance
//     * FROM records
//     * WHERE words CONTAINS ?
//     * ALLOW FILTERING
//     *
//     * @param words
//     * @return
//     */
//   def getFactsByUsedWord(word: String): List[Records] = {
//     run {
//       query[Records]
//         .filter(r => r.words.contains(lift(word)))
//         .allowFiltering
//     }
//   }

//   /** incrementFactToWordImportance
//     *
//     * UPDATE records
//     * SET importance = importance + 1
//     * WHERE record_uuid = ?
//     *
//     * @param id
//     */
//   def incrementFactToWordImportance(record_uuid: UUID) = {
//     run {
//       query[Records]
//         .filter(_.record_uuid == lift(record_uuid))
//         .update(p => p.importance -> (p.importance + 1))
//     }
//     true
//   }

//   /** getFactsByUuids
//     *
//     * SELECT record_uuid, account_uuid, tags, words, record, importance
//     * FROM records
//     * WHERE record_uuid IN (?)
//     *
//     * @param factIds
//     * @return
//     */
//   def getFactsByUuids(factIds: List[UUID]): List[Records] = {
//     run(query[Records].filter(a => liftQuery(factIds).contains(a.record_uuid)))
//   }

//   /** upsertFact
//     *
//     * UPDATE records
//     * SET
//     *   record_uuid = ?,
//     *   account_uuid = ?,
//     *   tags = ?,
//     *   words = ?,
//     *   record = ?,
//     *   importance = ?
//     * WHERE record_uuid = ?
//     *   AND account_uuid = ?
//     *
//     * @param fact
//     * @return
//     */
//   def upsertFact(fact: Records): Boolean = {
//     run {
//       query[Records]
//         .filter(a =>
//           (a.record_uuid == lift(fact.record_uuid))
//             && a.account_uuid == lift(fact.account_uuid)
//         )
//         .update(
//           _.record_uuid -> lift(fact.record_uuid),
//           _.account_uuid -> lift(fact.account_uuid),
//           _.tags -> lift(fact.tags),
//           _.words -> lift(fact.words),
//           _.record -> lift(fact.record),
//           _.importance -> lift(fact.importance)
//         )
//     }
//     true
//   }

//   /** batchInsertWordsToFact
//     *
//     * INSERT INTO records (
//     *   record_uuid,
//     *   account_uuid,
//     *   tags,
//     *   words,
//     *   record,
//     *   importance) 
//     * VALUES (?, ?, ?, ?, ?, ?)
//     *
//     * @param insertValues
//     */
//   def batchInsertWordsToFact(insertValues: List[Records]) = {
//     run {
//       quote {
//         liftQuery(insertValues).foreach(fact =>
//           query[Records]
//             .insert(
//               _.record_uuid -> lift(fact.record_uuid),
//               _.account_uuid -> lift(fact.account_uuid),
//               _.tags -> lift(fact.tags),
//               _.words -> lift(fact.words),
//               _.record -> lift(fact.record),
//               _.importance -> lift(fact.importance)
//             )
//         )
//       }
//     }
//     true
//   }

//   /** getIdsForWords
//     *
//     * SELECT record_uuid, account_uuid, tags, words, record, importance 
//     * FROM records 
//     * WHERE words IN (?)
//     * ALLOW FILTERING
//     *
//     * @param words
//     */
//   def getIdsForWords(words: List[String]): List[Records] = {
//     run {
//       query[Records]
//         .filter(a => liftQuery(words).contains(a.words))
//         .allowFiltering
//     }
//   }
// }
