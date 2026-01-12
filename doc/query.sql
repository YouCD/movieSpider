select  from_unixtime( timestamp),web from feed_video where web in (select distinct feed_video.web from feed_video) ;



SELECT from_unixtime(fv1.timestamp), fv1.web
FROM feed_video fv1
WHERE NOT EXISTS (
    SELECT 1
    FROM feed_video fv2
    WHERE fv2.web = fv1.web
      AND fv2.timestamp > fv1.timestamp
);
SELECT from_unixtime(timestamp) as time, web
FROM (
    SELECT timestamp, web,
    ROW_NUMBER() OVER (PARTITION BY web ORDER BY timestamp DESC) as rn
    FROM feed_video
    ) t
WHERE rn = 1;


SELECT from_unixtime(fv.timestamp), fv.web
FROM feed_video fv
         INNER JOIN (
    SELECT web, MAX(timestamp) as max_timestamp
    FROM feed_video
    GROUP BY web
) latest ON fv.web = latest.web AND fv.timestamp = latest.max_timestamp;



select  from_unixtime( timestamp),web from feed_video where web="tpbpirateproxy" order by timestamp desc limit 1;





update feed_video set web="other" where web="TgxDump";

# 获取磁力链接的种子
select * from feed_video where web="tpbpirateproxy" and torrent_url!="" order by timestamp desc limit 1;


#  查询每个站点最新的时间
select *
from (SELECT from_unixtime(fv.timestamp), fv.web
      FROM feed_video fv
               INNER JOIN (SELECT web, MAX(timestamp) as max_timestamp
                           FROM feed_video
                           GROUP BY web) latest ON fv.web = latest.web AND fv.timestamp = latest.max_timestamp) as t
group by web
order by t.`from_unixtime(fv.timestamp)` desc;

# 统计指定时间点 每个站点的种子数
select count(*) as cnt, web from feed_video  where timestamp>=1767801600 group by web  ;